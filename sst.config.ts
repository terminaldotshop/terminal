/// <reference path="./.sst/platform/config.d.ts" />

import { resolve } from "path";

export default $config({
  app(input) {
    return {
      name: "terminal-shop",
      removal: input?.stage === "production" ? "retain" : "remove",
      home: "cloudflare",
      providers: {
        aws: true,
        random: true,
        docker: true,
      },
    };
  },
  async run() {
    const domain =
      $app.stage === "production"
        ? "terminal.shop"
        : $app.stage + ".dev.terminal.shop";
    const secrets = {
      SwellSecret: new sst.Secret("SwellSecret"),
    };
    const auth = new sst.cloudflare.Auth("Auth", {
      authenticator: {
        link: [secrets.SwellSecret],
        handler: "./packages/workers/src/auth.ts",
        domain: "auth." + domain,
      },
    });
    const api = new sst.cloudflare.Worker("Api", {
      handler: "./packages/workers/src/api.ts",
      link: [secrets.SwellSecret, auth],
      domain: "api." + domain,
    });
    if (!$dev) {
      const www = new sst.cloudflare.StaticSite("Www", {
        domain: "www." + domain,
        path: "./packages/www",
        build: {
          command: "bun run build",
          output: "./dist",
        },
      });

      const repository = new aws.ecr.Repository("DockerRepository", {
        name: [$app.name, $app.stage].join("-"),
        forceDelete: true,
      });

      const vpc = new aws.ec2.Vpc("Vpc", {
        cidrBlock: "10.0.0.0/16",
        enableDnsSupport: true,
        enableDnsHostnames: true,
      });
      const subnet = new aws.ec2.Subnet("VpcSubnet", {
        vpcId: vpc.id,
        cidrBlock: "10.0.1.0/24",
        mapPublicIpOnLaunch: true,
      });
      const igw = new aws.ec2.InternetGateway("VpcIgw", {
        vpcId: vpc.id,
      });
      const routeTable = new aws.ec2.RouteTable("VpcRouteTable", {
        vpcId: vpc.id,
        routes: [
          {
            cidrBlock: "0.0.0.0/0",
            gatewayId: igw.id,
          },
        ],
      });
      const registryInfo = repository.registryId.apply(async (registryId) => {
        const credentials = await aws.ecr.getCredentials({
          registryId: registryId,
        });
        const decodedCredentials = Buffer.from(
          credentials.authorizationToken,
          "base64",
        ).toString();
        const [username, password] = decodedCredentials.split(":");
        return {
          server: credentials.proxyEndpoint,
          username: username,
          password: password,
        };
      });

      const image = new docker.Image("SSHImage", {
        build: {
          context: resolve("./go"),
          platform: "linux/amd64",
        },
        imageName: $interpolate`${repository.repositoryUrl}:${$app.stage}`,
        registry: registryInfo,
      });
      new aws.ec2.RouteTableAssociation("VpcRouteTableAssociation", {
        subnetId: subnet.id,
        routeTableId: routeTable.id,
      });
      const cluster = new aws.ecs.Cluster("Cluster");
      const executionRole = new aws.iam.Role("SSHRole", {
        assumeRolePolicy: aws.iam.assumeRolePolicyForPrincipal({
          Service: "ecs-tasks.amazonaws.com",
        }),
      });
      new aws.iam.RolePolicyAttachment("ExecutionRolePolicyAttachment", {
        role: executionRole.name,
        policyArn:
          "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy",
      });
      new aws.iam.RolePolicyAttachment("SSHRolePolicyAttachment", {
        role: executionRole.name,
        policyArn:
          "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy",
      });

      const portSSH = 22;
      const portHTTP = 80;

      const taskDefinition = new aws.ecs.TaskDefinition("SSHTask", {
        family: "ssh",
        trackLatest: true,
        cpu: "2048",
        memory: "4096",
        networkMode: "awsvpc",
        requiresCompatibilities: ["FARGATE"],
        executionRoleArn: executionRole.arn,
        containerDefinitions: $jsonStringify([
          {
            name: "ssh",
            image: image.repoDigest,
            portMappings: [
              {
                containerPort: portSSH,
                hostPort: portSSH,
                protocol: "tcp",
              },
              {
                containerPort: portHTTP,
                hostPort: portHTTP,
                protocol: "tcp",
              },
            ],
            environment: [
              {
                name: "SSH_PORT",
                value: portSSH.toString(),
              },
              {
                name: "HTTP_PORT",
                value: portHTTP.toString(),
              },
            ],
          },
        ]),
      });
      const sshSecurityGroup = new aws.ec2.SecurityGroup(
        "SSHNlbSecurityGroup",
        {
          vpcId: vpc.id,
          egress: [
            {
              fromPort: 0,
              toPort: 0,
              protocol: "-1",
              cidrBlocks: ["0.0.0.0/0"],
            },
          ],
          ingress: [
            {
              fromPort: portSSH,
              toPort: portSSH,
              protocol: "tcp",
              cidrBlocks: ["0.0.0.0/0"],
            },
            {
              fromPort: portHTTP,
              toPort: portHTTP,
              protocol: "tcp",
              cidrBlocks: ["0.0.0.0/0"],
            },
          ],
        },
      );

      const sshTargetGroup = new aws.lb.TargetGroup("SSHNlbTargetGroup", {
        port: portSSH,
        protocol: "TCP",
        targetType: "ip",
        vpcId: vpc.id,
      });
      const httpTargetGroup = new aws.lb.TargetGroup("NlbTargetGroupHttp", {
        port: portHTTP,
        protocol: "TCP",
        targetType: "ip",
        vpcId: vpc.id,
      });
      const service = new aws.ecs.Service("SSHService", {
        cluster: cluster.arn,
        taskDefinition: taskDefinition.arn,
        desiredCount: 1,
        launchType: "FARGATE",
        networkConfiguration: {
          assignPublicIp: true,
          subnets: [subnet.id],
          securityGroups: [sshSecurityGroup.id],
        },
        loadBalancers: [
          {
            targetGroupArn: sshTargetGroup.arn,
            containerName: "ssh",
            containerPort: portSSH,
          },
          {
            targetGroupArn: httpTargetGroup.arn,
            containerName: "ssh",
            containerPort: portHTTP,
          },
        ],
      });

      const nlb = new aws.lb.LoadBalancer("SSHNlb", {
        internal: false,
        loadBalancerType: "network",
        subnets: [subnet.id],
        enableCrossZoneLoadBalancing: true,
        securityGroups: [sshSecurityGroup.id],
      });
      new aws.lb.Listener("SSHListener", {
        loadBalancerArn: nlb.arn,
        port: portSSH,
        protocol: "TCP",
        defaultActions: [
          {
            type: "forward",
            targetGroupArn: sshTargetGroup.arn,
          },
        ],
      });
      new aws.lb.Listener("HttpListener", {
        loadBalancerArn: nlb.arn,
        port: portHTTP,
        protocol: "TCP",
        defaultActions: [
          {
            type: "forward",
            targetGroupArn: httpTargetGroup.arn,
          },
        ],
      });
    }
    return {
      api: api.url,
      auth: auth.url,
    };
  },
});

/// <reference path="./.sst/platform/config.d.ts" />
import { resolve } from "path";
export default $config({
  app(input) {
    return {
      name: "terminal-shop",
      removal: input?.stage === "production" ? "retain" : "remove",
      home: "cloudflare",
      providers: {
        aws: {
          profile: process.env.GITHUB_ACTIONS
            ? undefined
            : input.stage === "production"
              ? "ironbay-production"
              : "ironbay-dev",
        },
        random: true,
        docker: true,
        tls: true,
      },
    };
  },
  async run() {
    const isPermanentStage =
      $app.stage === "production" || $app.stage === "dev";
    const domain =
      $app.stage === "production"
        ? "terminal.shop"
        : $app.stage + ".dev.terminal.shop";
    const secrets = {
      SwellSecret: new sst.Secret("SwellSecret"),
      AirtableSecret: new sst.Secret("AirtableSecret"),
      StripeSecret: new sst.Secret("StripeSecret"),
      ShippoSecret: new sst.Secret("ShippoSecret"),
    };
    const auth = new sst.cloudflare.Auth("Auth", {
      authenticator: {
        link: [secrets.SwellSecret, secrets.StripeSecret],
        handler: "./packages/workers/src/auth.ts",
        domain: "auth." + domain,
      },
    });
    const api = new sst.cloudflare.Worker("Api", {
      handler: "./packages/workers/src/api.ts",
      link: [
        secrets.SwellSecret,
        secrets.AirtableSecret,
        secrets.StripeSecret,
        secrets.ShippoSecret,
        auth,
      ],
      domain: "api." + domain,
    });
    const www = new sst.cloudflare.StaticSite("Www", {
      domain: "www." + domain,
      path: "./packages/www",
      environment: {
        PUBLIC_API_URL: api.url.apply((u) => u!),
      },
      build: {
        command: "bun run build",
        output: "./dist",
      },
    });
    if (isPermanentStage) {
      const github = new aws.iam.OpenIdConnectProvider("GithubOidc", {
        url: "https://token.actions.githubusercontent.com",
        clientIdLists: ["sts.amazonaws.com"],
        thumbprintLists: [
          "6938fd4d98bab03faadb97b34396831e3780aea1",
          "1c58a3a8518e8759bf075b76b750d4f2df264fcd",
        ],
      });
      const githubRole = new aws.iam.Role("GithubRole", {
        name: [$app.name, $app.stage, "github"].join("-"),
        assumeRolePolicy: {
          Version: "2012-10-17",
          Statement: [
            {
              Effect: "Allow",
              Principal: {
                Federated: github.arn,
              },
              Action: "sts:AssumeRoleWithWebIdentity",
              Condition: {
                StringLike: github.url.apply((url) => ({
                  [`${url}:sub`]: "repo:terminalhq/terminal:*",
                })),
              },
            },
          ],
        },
      });
      new aws.iam.RolePolicyAttachment("GithubRolePolicy", {
        policyArn: "arn:aws:iam::aws:policy/AdministratorAccess",
        role: githubRole.name,
      });
    }
    if (!$dev) {
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
      const portHTTPS = 443;
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
            {
              fromPort: portHTTPS,
              toPort: portHTTPS,
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
      const cert = new aws.acm.Certificate("SSLCertificate", {
        domainName: "terminal.shop",
        validationMethod: "DNS",
      });
      const zone = await cloudflare.getZone({ name: "terminal.shop" });
      const records: cloudflare.Record[] = [];
      cert.domainValidationOptions.apply((domainValidationOptions) => {
        const [options] = domainValidationOptions;
        records.push(
          new cloudflare.Record("CertificateValidationRecord", {
            zoneId: zone.zoneId,
            name: options.resourceRecordName,
            type: options.resourceRecordType,
            value: options.resourceRecordValue,
            ttl: 300,
          }),
        );
      });
      const validation = new aws.acm.CertificateValidation("CertValidation", {
        certificateArn: cert.arn,
        validationRecordFqdns: records.map((record) => record.hostname),
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
      new aws.lb.Listener("HttpsListener", {
        certificateArn: validation.certificateArn,
        loadBalancerArn: nlb.arn,
        port: portHTTPS,
        protocol: "TLS",
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

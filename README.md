### Setup

- install bun https://bun.sh/
- install sst https://sst.dev/
- go [here](https://dash.cloudflare.com/profile/api-tokens?permissionGroupKeys=%5B%7B%22key%22%3A%22account_settings%22%2C%22type%22%3A%22edit%22%7D%2C%7B%22key%22%3A%22dns%22%2C%22type%22%3A%22edit%22%7D%2C%7B%22key%22%3A%22memberships%22%2C%22type%22%3A%22read%22%7D%2C%7B%22key%22%3A%22user_details%22%2C%22type%22%3A%22edit%22%7D%2C%7B%22key%22%3A%22workers_kv_storage%22%2C%22type%22%3A%22edit%22%7D%2C%7B%22key%22%3A%22workers_r2%22%2C%22type%22%3A%22edit%22%7D%2C%7B%22key%22%3A%22workers_routes%22%2C%22type%22%3A%22edit%22%7D%2C%7B%22key%22%3A%22workers_scripts%22%2C%22type%22%3A%22edit%22%7D%2C%7B%22key%22%3A%22workers_tail%22%2C%22type%22%3A%22read%22%7D%5D&name=sst&accountId=*&zoneId=all)
- make sure Account Resources is limited to `terminal` instead of `All accounts` and create the api token.
- save it to a `.env` file like this
```
CLOUDFLARE_API_TOKEN=xxx
```
- place this in `~/.aws/config`
```
[sso-session terminal]
sso_start_url = https://terminaldotshop.awsapps.com/start
sso_region = us-east-2

[profile terminal-dev]
sso_session = terminal
sso_account_id = 058264103289
sso_role_name = AdministratorAccess
region = us-east-2

[profile terminal-production]
sso_session = terminal
sso_account_id = 211125775473
sso_role_name = AdministratorAccess
region = us-east-2
```
- need to login once a day with `bun sso` in root

### Frontend only

if you're only working on frontend you don't need to bring up the whole system - you can create a shell that's linked to the `dev` environment.

- `cd go`
- `sst shell --stage=dev` - this connects you to the dev environment and opens a bash shell
- `go run ./cmd/cli` - this will run the cli

note - this loads dev secrets into your environment. most of them aren't them sensitive but if you are streaming you should avoid logging your env.


### Full system

if you're working on the full system you can do

- `sst dev` in root to bring up everything


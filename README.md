# discord weather notification

This project is a serverless application built with AWS Lambda, EventBridge. It retrieves weather forecasts each morning at 7:00 AM JST and sends notifications to Discord via webhook if rain is forecast.

## Environment Variables & Secrets Management

- DISCORD_WEBHOOK_URL
  - This is a sensitive piece of information.
  - Use a Terraform variable (defined in a .tfvars file)

## Using make.sh

A `make.sh` script is provided to automate the build, package, and deploy processes. It performs:

Lambda function build and packaging.

Terraform initialization, planning, and apply (with a manual confirmation step for the plan).

Before running, ensure it has execution permission:

```bash
chmod +x make.sh
./make.sh
```

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.

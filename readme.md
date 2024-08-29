# Project Setup and Deployment Guide

This guide provides instructions on how to set up and deploy the application using Docker. 

## Prerequisites

Make sure the following packages are installed on your system:

1. **Docker**: [Install Docker](https://docs.docker.com/engine/install/)
2. **yq**: Install `yq` with the following commands:
    ```bash
    sudo wget https://github.com/mikefarah/yq/releases/latest/download/yq_linux_amd64 -O /usr/bin/yq
    sudo chmod +x /usr/bin/yq
    ```
3. **Check Docker Compose Version**: 
    - If you are using Docker Compose version 1, the execution script must be `docker-compose`.
    - If you are using Docker Compose version 2, the execution script must be `docker compose`.
    - Make sure your version's execution script matches the one used in `./starter.sh`. If not, update the script by replacing the execution command on lines 147, 158, and 168.

## First-Time Setup

If you are setting up the application for the first time, follow these steps:

1. **Create Docker Network**: Ensure that the Docker network is created. The network name must match the one specified in the `docker-compose` file.
   
2. **Run the Starter Script**:
    - Execute the `starter.sh` script by running:
      ```bash
      ./starter.sh --env--
      ```
    - This script will automatically create the database, Docker volumes, and various configurations using Docker. The application will be ready to use after running this script.
    - Replace `--env--` with the environment you want to use. There are two environments available: `dev` and `production`. The environment you choose will determine which files are executed:
        - **Production**: Executes `env.production`, `Dockerfile`, and `docker-compose.production.yml`.
        - **Dev**: Executes `env.dev`, `Dockerfile`, and `docker-compose.dev.yml`.

## Redeploying the Application

If you've previously installed the application and plan to deploy the latest code:

1. **Update Docker Image and Container Names**:
    - Update the image and container names in the `docker-compose` file. This prevents using the same image or container, which might cause the code changes to be ignored after deployment.
    - You can use versioning for the names, for example:
      ```yaml
      image: turbine-api-production:v1.0.0
      container_name: turbine-api-production-v100
      ```
    - The previous container will automatically be replaced with the new one, and the latest application version will be ready for use.

---

Feel free to reach out if you encounter any issues during the setup or deployment process.

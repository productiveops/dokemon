# Dokémon

Dokémon is a friendly GUI for managing Docker Containers on Virtual Machines.

_Website URL:_ https://dokemon.dev

## Quickstart

You can run the below commands to quickly try out Dokémon.

**Note:** Whenever possible, it is recommended that you run Dokémon in a private network and do not expose it to the Internet. In cases where this is not possible, for example when running on a VPS to which you only have public access, you should run Dokémon behind an SSL enabled proxy and use a strong password for maximum security. Refer the next section for sample configuration using Traefik.

    # Create directory to store Dokemon data
    mkdir ./dokemondata

    # Run Dokemon
    sudo docker run -p 9090:9090 \
                -v ./dokemondata:/data \
                -v /var/run/docker.sock:/var/run/docker.sock \
                --restart unless-stopped \
                --name dokemon -d productiveops/dokemon:latest

## Using Traefik with LetsEncrypt SSL certificate

This is an example configuration for running Dokémon behind Traefik with LetsEncrypt SSL certificate.

**Note:** This is a sample configuration. Please modify it as per your requirements.

    version: "3.3"

    services:
      traefik:
        image: "traefik:v2.10"
        container_name: "traefik"
        command:
          - "--log.level=DEBUG"
          - "--accesslog=true"
          - "--api.insecure=true"
          - "--providers.docker=true"
          - "--providers.docker.exposedbydefault=false"
          - "--entrypoints.websecure.address=:443"
          - "--certificatesresolvers.dokemon.acme.tlschallenge=true"
          - "--certificatesresolvers.dokemon.acme.email=your.email@example.com"
          - "--certificatesresolvers.dokemon.acme.storage=/letsencrypt/dokemon.json"
        ports:
          - "443:443"
          - "8080:8080"
        volumes:
          - "./letsencrypt:/letsencrypt"
          - "/var/run/docker.sock:/var/run/docker.sock:ro"

      dokemon:
        image: productiveops/dokemon:latest
        container_name: dokemon
        restart: unless-stopped
        labels:
          - "traefik.enable=true"
          - "traefik.http.routers.dokemon.rule=Host(`dokemon.example.com`)"
          - "traefik.http.routers.dokemon.entrypoints=websecure"
          - "traefik.http.routers.dokemon.tls.certresolver=dokemon"
        ports:
          - 9090:9090
        volumes:
          - ./dokemondata:/data
          - /var/run/docker.sock:/var/run/docker.sock

1. Create a file name `compose.yaml` on your server. Copy and paste the above YAML definition into the file. Modify the email and host. Make any other changes required as per your requirements.
2. Run `mkdir ./letsencrypt && mkdir ./dokemondata`
3. Run `docker compose up -d`

Open https://dokemon.example.com (your URL here) in the browser. It can take a few seconds for the SSL certificate to be provisioned. If you get SSL error, please wait for a few moments and refresh your browser.

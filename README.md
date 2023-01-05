# dvb (Docker Volume Backup)

- [dvb (Docker Volume Backup)](#dvb-docker-volume-backup)
  - [Requirements](#requirements)
  - [Install](#install)
  - [Config](#config)
    - [Variables](#variables)
    - [Backup](#backup)
      - [Docker](#docker)
    - [Destination](#destination)
    - [Retain](#retain)
      - [Local](#local)
    - [Alerts](#alerts)
      - [Discord Webhooks](#discord-webhooks)
      - [Email](#email)
    - [Daemon](#daemon)
  - [Full Config Example](#full-config-example)

Docker Volume Backup is a cli tool I use to backup my data on my docker servers.  I have mixed luck with cron trying to chain together commands to get them to work the way I want.  I want also wnt to see the notifications come in when the backups are done so I can keep an eye on things.

## Requirements

1. Have Go 1.19 installed or newer (no plans to release compiled versions at this time)
2. Have Docker CLI installed (This was built against version 20.10.14)

## Install

`go install github.com/jtom38/dvb@latest`

## Config

The application does require a config file to work correctly.  The config is too complex to run inline with the CLI

### Variables

The config file does support a small subset of variables.

- `{{PWD}}` = Current working directory
- `{{DATE}}` = Todays date in the following format `20221201`
- `{{USERDIR}}` = This will find the running users directory

If you use one of these you can pre/post append values and things will be replaced and updated.

Example: {{USERDIR}}/nas = /home/username/nas (on linux)

### Backup

The `Backup` flag is a top level object that defines all the target that the tool will backup.

#### Docker

This uses the Docker CLI tool to backup your containers.  I have not tested podmon to validate if this works out the box or not.

- Name string: defines the name of the container to target.
- Directory string defines where inside the container to target to backup data.
- Tar.Directory string: Defines where the backup file will be created.
- Tar.Pattern string: The backup file name pattern.  .tar is appended to the file.
- Post.Reboot: array string - optional: Defines any extra containers that should be rebooted after the backup has been performed.  This can be used to make sure any dependant apps can come back in a clean state if you take its database offline for example.

```yaml
Backup:
  Docker:
    - Name: app-db-01
      Directory: /var/lib/data
      Tar:
        Directory: "{{PWD}}"
        Pattern: "app-db-{{DATE}}"

      Post:
        Reboot:
          - app-api-01
```

### Destination

This tells the app what to do with the backups once they have been made.  Right now, it only supports moving data around on your own host.

### Retain

This block tells the app how many backups you want to retain on disk.  When this is present the app will check the destination before any actions are taken to make sure it can store the data as expected.  If a file already exists with the same name, it will append a .1 (or higher) till it finds a name that isn't taken.

If the app runs into any extra files, it will remove the oldest backup in the folder.  Make sure that the user running DVB will be able to read, write and delete out of the folder if you use the retain statement.  

This is an optional part of the config, if you don't want it, comment/delete it from your config.

```yaml
Backup:
  ...

Destination:
  Retain:
    Days: 10

Alert:
  ...
```

#### Local

The Local statement will move the generated tiles to a different location on the same device.  If you have a SMB or NFS mount on your system bound to a directory, then DVB will move the data to that folder.

This will create a subfolder with the containers name and load the backups into that directory.  The config below will create a directory `/mnt/nas/backups/webdav` and move the backups to that directory for you.

```yaml
Backup:
  Docker:
    - Name: webdav
      ...

Destination:
  Local:
    Path: /mnt/nas/backups
```

### Alerts

DVB will send out alerts with a log dump if you want them.

Alerts is defined as a top level object in the config and you can define the following attributes.

#### Discord Webhooks

- `Username` = The name that will be used when message is set
- `Webhooks` = This contains all the webhooks that you want to send to.  You can send to multiple if you want.

```yaml
Backup:
  ...

Destination:
  ...

Alert:
  Discord:
    UserName: DVB02
    Webhooks:
      - https://discord.com/api/webhooks...
      - https://discord.com/api/webhooks...
        
```

#### Email

To send emails with DVB you will need to make sure that the account you are using is able to send via SMTP.  If you use gmail, make sure you have two factor setup and generate an app password.  Use that app password with this application and you should see emails trigger.

I am giving the gmail config as its what I have tested with.  Any email account with SMTP should work with this.

```yaml
Backup:
  ...

Destination:
  ...

Alert:
  Email:
    Account:
      Username: serviceaccount@gmail.com
      Password: ThisIsNotARealPassword
      Host: smtp.gmail.com
      Port: 587
      UseTls: true
    From: serviceaccount@gmail.com
    To: servermaintainer@gmail.com
```

### Daemon

The daemon mode will keep the process alive and will generate backups based on a cron value defined in the yaml.

If you need assistance to build a cron timer, use [crontab.guru](https://crontab.guru).

```yaml

Daemon:
  # This will start the backup job at midnight 
  Cron: "0 0 * * *"

Backup:
  ...

Destination:
  ...

Alert:
  ...
```

## Full Config Example

```yaml
Daemon:
  Cron: "32 11 * * *"

Backup:
  Docker:
    - Name: webdav-app-1
      Directory: /var/lib/dav
      Tar:
        Directory: $PWD
        Pattern: webdav-data-{{date}}
        Extension: tar
      Post:
        Reboot:
          - webdav2
          - webdav3

Destination:
  Retain:
    Days: 10

  Local: 
    Path: /mnt/nas/backups

Alert:
  Discord:
    Username: docker-server-backups
    Webhooks:
      - https://discord.com/api/webhooks/...
  Email:
    Account:
      Username: serviceaccount@gmail.com
      Password: ThisIsNotARealPassword
      Host: smtp.gmail.com
      Port: 587
      UseTls: true
    From: serviceaccount@gmail.com
    To: servermaintainer@gmail.com
```

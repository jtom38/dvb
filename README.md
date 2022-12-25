# dvb (Docker Volume Backup)

Docker Volume Backup is a cli tool I use to backup my data on my docker servers.  I have mixed luck with cron trying to chain together commands to get them to work the way I want.  I want also wnt to see the notifications come in when the backups are done so I can keep an eye on things.

## Requirements

1. Have Go 1.19 installed or newer (no plans to release compiled versions at this time)
2. Have Docker CLI installed (This was built against version 20.10.14)

## Install

`go install github.com/jtom38/dvb@latest`

## Config

```yaml
Backup:
  # This defines that the containers come from Docker
  Docker:
      # Name of the container
    - Name: webdav-app-1
      # Where inside your container is the data that needs to be backed up
      Directory: /var/lib/dav
      Tar:
        # if $PWD is used, we convert this to your working directory as you expect.
        Directory: $PWD
        # {{date}} is replaced with a date code value 20221217
        Pattern: webdav-data-{{date}}
        Extension: tar
  
# Once its backed up, what do we do with the file?
# If you want to leave it in the current directory it was made, you don't need this block
Destination:
  Retain:
    # This tells us to remove files older then 10 days
    # if you want to retain everything, set this to 0
    Days: 10

  # Local flag defines that this will move data within the same host
  Local: 
    # Defines what directory that it will move to.
    # The container name will be append to the path
    Path: /mnt/nas/backups

# If you want this to send you alerts, right now discord is only supported
# If you don't want discord alerts, remove this block
Alert:
  Discord:
    Username: docker-server-backups
    Webhooks:
      - https://discord.com/api/webhooks/...

```

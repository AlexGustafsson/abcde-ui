# abcde-ui

A dead simple server/UI for the abcde ripper. Useful for running the tool in
Docker or Kubernetes.

The happy path is when you want a mostly manual ripping process, where you
manually change CDs, press the button in the UI and then manually curate the
metadata using something like
[MusicBrainz' Picard](https://picard.musicbrainz.org).

As the server uses abcde, a bad rip can be retried, or stopped and picked up at
a later time. If an issue arises needing more debugging, you can attach a
terminal to the container and investigate manually using abcde, flac or
libcdio-utils.

The server be run fully rootless and without any privileges when combined with
[ARM's smarter device plugin](https://gitlab.com/arm-research/smarter/smarter-device-manager).

Example Kubernetes deployment below.

```yaml
---
kind: Deployment
apiVersion: apps/v1
metadata:
  name: abcde
  labels:
    app: abcde
spec:
  selector:
    matchLabels:
      app: abcde
  strategy:
    type: Recreate
  template:
    metadata:
      labels:
        app: abcde
    spec:
      securityContext:
        runAsNonRoot: true
        runAsUser: 2000
        runAsGroup: 2000
        fsGroup: 2000
      containers:
        - name: abcde
          image: abcde-ui
          imagePullPolicy: Always
          securityContext:
            allowPrivilegeEscalation: false
            readOnlyRootFilesystem: true
            privileged: false
            capabilities:
              drop:
                - ALL
          workingDir: /var/data
          volumeMounts:
            - name: data
              mountPath: /var/data
            - name: tmp
              mountPath: /tmp
          ports:
            - name: web
              containerPort: 8080
          resources:
            requests:
              smarter-devices/sr0: 1
            limits:
              smarter-devices/sr0: 1
      volumes:
        - name: data
          hostPath:
            path: /home/example/media
        - name: tmp
          emptyDir: {}
---
# ...service
# ...ingress(es)
# ...certificate
```

## Tips

### Manually ripping CDs

If a CD rip fails to start, typically due to failed lookups of singles, you can
attach a terminal and run a command like this, which will ignore lookups and
just rip the cd.

```shell
abcde -n -Q cddb -N -o flac -p
```

Remember that in this case, abcde-ui can't create a unique directory for the rip
so be sure to rename the output directory so that it doesn't get overwritten.

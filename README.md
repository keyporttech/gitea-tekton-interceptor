# gitea-tektoncd-event interceptor

This a custom gitea interceptor for [tekton triggers](https://github.com/tektoncd/triggers). It is used by [k8sCI](https://github.com/keyporttech/k8sci) to validate gitea webhook requests.

# usage

This is for use as a webhook interceptor so that a tekton event listener can process gitea webhooks. See [tekton triggers event listeners](https://github.com/tektoncd/triggers/blob/master/docs/eventlisteners.md) for more details.

The docker image is published at [keyporttech/gitea-tektconcd-event-interceptor] (https://hub.docker.com/repository/docker/keyporttech/gitea-tektconcd-event-interceptor), and this image is used as service deployment in the [k8sCI](https://github.com/keyporttech/k8sci).

# building

### Prerequisites: golang, docker, makefile installed

```bash
make build #build
make docker # build docker image

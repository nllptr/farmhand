# Welcome to Farmhand!

Farmhand is a web service that enables a hobby farmer (such as myself) to better keep track of farm tasks that need to be done. Currently, I'm working to get the following pieces in place:

 - Weather reminders - if it's about get really cold the user should be reminded some time before, so that crops can be covered and animals put into proper shelter.

 - Livestock diseases - if diseases have been reported in the area, the user should be notified of this, so that their own animals can be put into quarantine.

 ## However..

 This is also very much a practice project for me to get some training with Kubernetes and related technologies. I don't expect this to ever be a "finished product".


## Running locally on a Kind cluster

1. Run `./helper.sh pre` to make sure your environment fulfills the prerequisites.
2. Run `./helper.sh setup` to create the cluster.
3. Run `skaffold run` to build the images.
4. Run `skaffold dev` to deploy all services to the cluster.

# TODO
- Generate skeleton .env file
- Convert the helper script to a cross-platform go program.
- The helper script has a lot of duplicated code. Make it nicer.

# welcomer
The welcomer is a microservice written in golang. It's a part of my experimental effort to implement Open Tracing using linkerd.
It calls the microservice [guesttracker](https://github.com/niksw7/guesttracker). 
How to view the logs.
brew install stern and then the following command
stern welcomer --container welcomer --all-namespaces


#how to upgrade
docker build -t=welcomer:1.x <br/>
change deployment.yaml to point to new version of welcomer 1.x<br/>
kubectl delete deployments nginx-ingress -n hackerspace
#!/bin/bash
exitIfReturnCodeNonZero(){ 
      returnvalue=$?
      if [ $returnvalue -ne 0 ]; then
            echo "Something went wrong in $1 ..Exiting"
            exit $returnvalue
      fi
}
# cleanup [start]
kubectl delete deployments guesttracker welcomer nginx-ingress -n hackerspace
kubectl delete deployments guesttracker welcomer nginx-ingress -n hackerspace
# cleanup [end]

currentpath=$PWD
# start welcomer and guesttracker [start]
cd $GOPATH/src/github.com/welcomer
kubectl apply -f services.yaml -f deployment.yaml
exitIfReturnCodeNonZero "starting welcomer"

cd $GOPATH/src/github.com/guesttracker
kubectl apply -f services.yaml -f deployment.yaml
exitIfReturnCodeNonZero "starting guesttracker"
# start welcomer and guesttracker [end]


#Inject linkerd plane proxies
kubectl get -n hackerspace deploy -o yaml \
  | linkerd inject - \
  | kubectl apply -f -
exitIfReturnCodeNonZero "injecting linkerd plane proxies"

#Start openconsensus collector and jaeger and ingress
cd $GOPATH/src/github.com/welcomer
linkerd inject tracing.yml | kubectl apply -f -
exitIfReturnCodeNonZero "starting openconsensus collector and jaeger"
kubectl apply -f ingress.yml
exitIfReturnCodeNonZero "starting ingress"

kubectl -n tracing port-forward deploy/jaeger 16686 & ; open http://localhost:16686 &

echo "Successfully deployed"
cd $currentpath

#Deleting linkerd pods
#kubectl delete deployments $( kubectl get deployments -n linkerd | cut -d ' ' -f1) -n linkerd

#curl -sL https://run.linkerd.io/install | sh
#export PATH=$PATH:$HOME/.linkerd2/bin


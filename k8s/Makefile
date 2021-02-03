install:
	helm upgrade -i prometheus-demo ./prometheus-demo -n prometheus-demo --create-namespace
port-forward:
	kubectl --namespace prometheus-demo port-forward svc/prometheus-demo 8080
configmap:
	kubectl create configmap makefile --from-file=../Makefile -n prometheus-demo
apply_load:
	kubectl apply -f load/ -n prometheus-demo
delete_load:
	kubectl delete -f load/ -n prometheus-demo


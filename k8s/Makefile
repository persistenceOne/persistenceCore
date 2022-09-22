PROCESS=persistencecore
NAMESPACE=local

apply:
	kubectl apply -k overlays/$(CLUSTER)/$(PROCESS)

delete:
	kubectl delete -k overlays/$(CLUSTER)/$(PROCESS)


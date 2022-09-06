PROCESS=persistencecore
CLUSTER=devnet

apply:
	kubectl apply -k overlays/$(CLUSTER)/$(PROCESS)

delete:
	kubectl delete -k overlays/$(CLUSTER)/$(PROCESS)


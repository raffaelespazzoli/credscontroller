FROM registry.access.redhat.com/rhel:latest
LABEL authors="Raffaele Spazzoli <rspazzol@redhat.com>" 
ADD credscontroller /credscontroller
ENTRYPOINT ["/credscontroller"]

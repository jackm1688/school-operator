FROM registry.access.redhat.com/ubi8/ubi-minimal:latest

ENV OPERATOR=/usr/local/bin/school-operator \
    USER_UID=1001 \
    USER_NAME=school-operator \
    GO111MODEULE=on \
	CGO_ENABLED=0 \
	GOOS=linux \
	GOARCH=amd64 \
	GOPROXY=https://goproxy.cn,direct

# install operator binary
COPY build/_output/bin/school-operator ${OPERATOR}

COPY build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}

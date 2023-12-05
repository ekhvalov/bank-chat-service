FROM antonboom/writing-go-service.afc as builder

FROM alpine
COPY --from=builder /go/bin/afc /usr/local/bin/afc

CMD /usr/local/bin/afc


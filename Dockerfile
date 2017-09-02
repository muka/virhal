FROM scratch
COPY ./virhal /virhal
ENTRYPOINT ["/virhal"]

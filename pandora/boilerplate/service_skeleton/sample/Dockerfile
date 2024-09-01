FROM public.ecr.aws/o3d4h0u7/golang:1.17.1-alpine
WORKDIR /dist
COPY japi-session-service ./
RUN chmod +x japi-session-service
EXPOSE 3000
ENTRYPOINT ["./sample","server"]
#CMD ["cron", "fruadcheckcorrection"]

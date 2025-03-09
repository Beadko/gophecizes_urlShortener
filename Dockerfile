FROM cockroachdb/cockroach:v23.1.11

EXPOSE 26258 8081

VOLUME /cockroach/cockroach-data

CMD ["start-single-node", "--insecure"]


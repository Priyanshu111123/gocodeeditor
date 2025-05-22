# Use an official Java base image
FROM openjdk:17-slim

# Set working directory
WORKDIR /workspace

# Copy compile script
COPY compile.sh /usr/local/bin/compile.sh
RUN chmod +x /usr/local/bin/compile.sh

CMD ["/usr/local/bin/compile.sh"]

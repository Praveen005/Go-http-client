Why do we need to close the response body?

- Closing the response body allows the underlying resources (such as network connections) to be released promptly. Failing to close the body can lead to resource leaks, especially in scenarios with a high volume of HTTP requests.
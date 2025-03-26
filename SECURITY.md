## Security

> **Fairway Bridge is a proof-of-concept integration tool and is not considered "production ready." Do not expect robust security measures or enterprise-level protections.**

Fairway Bridge currently uses basic TCP data connections for supported launch monitors and simulators, and the HTTP server runs without HTTPS encryption.

It should not be deployed as a public-facing API or used in critical, security-sensitive environments without additional safeguards.

If you discover a new security vulnerability, please report it using the secure channels detailed below.

Additional security concerns may be identified as the project evolves.

---

## Reporting Security Issues

**Do not report security vulnerabilities via public GitHub issues.**

Please email your report to [git@stephenmendez.dev](mailto:git@stephenmendez.dev). You can expect a response within 72 hours.

Include as much of the following information as possible:

- **Type of Issue:** e.g., buffer overflow, SQL injection, cross-site scripting, etc.
- **Affected Files:** Full paths of the source file(s) where the issue manifests.
- **Code Location:** The tag/branch/commit or direct URL where the issue occurs.
- **Reproduction Steps:** Any special configuration and step-by-step instructions to reproduce the vulnerability.
- **Proof-of-Concept:** Exploit code or demonstrations of the issue (if possible).
- **Impact:** A detailed description of the potential impact, including how an attacker might exploit the vulnerability.

---

## Known Vulnerabilities

> In the interest of transparency, the following known vulnerabilities are documented as a reminder that Fairway Bridge should not be deployed as a public-facing API or used in critical, security-sensitive environments without additional safeguards.

### Server

- **Unencrypted Communication:**  
  The HTTP server transmits data in plaintext without encryption.  
  *Recommendation:* Run the application behind a secure reverse proxy (e.g., NGINX with HTTPS enabled).

- **Verbose Error Information:**  
  Exception details may be exposed in HTTP responses, potentially leaking sensitive information.

### Endpoints

- **Lack of Rate Limiting:**  
  No endpoint rate limiting is implemented, making the application susceptible to denial-of-service or brute-force attacks.

- **Input Validation:**  
  Endpoints currently do not enforce comprehensive input validation.  
  *Impact:* Malformed or malicious input could lead to unexpected behavior or data corruption.

### Authentication

- **No Authentication Framework:**  
  There is no authentication policy.

### Logging

- **Limited Logging Capabilities:**  
  Current logging mechanisms are detailed, and sensitive data may be logged.
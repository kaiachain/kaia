# Introduction

Thank you for helping keep the Kaia ecosystem secure.

We operate a responsible disclosure and bug bounty program in partnership with [HackenProof](https://hackenproof.com/companies/kaia). This document outlines how to report vulnerabilities, our bounty scope, and program rules.

---

##  Reporting a Vulnerability

Please **do not use GitHub, email, or Discord** to report vulnerabilities.

Instead, submit all reports via our official bug bounty dashboard:

👉 [Report a vulnerability on HackenProof](https://hackenproof.com/companies/kaia)

You must report vulnerabilities **within 24 hours of discovery** and exclusively via HackenProof to be eligible for a bounty.

---

##  Bug Bounty Program Overview

We offer bounties for valid, impactful vulnerabilities across the [Kaia Protocol](https://hackenproof.com/programs/kaia-protocol) and [Kaia Web](https://hackenproof.com/programs/kaia-web) ecosystem.

Reward amounts vary based on:
- Impact
- Severity
- Quality of report and PoC

---

## In-Scope Targets

###  Kaia Protocol (Blockchain Layer)

Focuses on blockchain protocol vulnerabilities, including but not limited to:
- Stealing or loss of funds
- Unauthorized or manipulated transactions
- Price or fee manipulation
- Balance or tokenomics manipulation
- Privacy violations
- Cryptographic flaws

###  Kaia Web (Web Apps, SDKs, APIs)

Focuses on web-based vulnerabilities such as:
- Business logic issues
- Payment manipulation
- Remote Code Execution (RCE)
- SQL/XXE Injection
- Access control issues (IDOR, Privilege Escalation)
- Sensitive data leaks
- SSRF, CSRF, XSS
- File inclusion, directory traversal

See the full list on our [Kaia Web HackenProof Program page](https://hackenproof.com/programs/kaia-web).

---

##  Out-of-Scope Vulnerabilities

Some issues are **not eligible** for bounties, including:

### Blockchain
- Network-level DoS
- Attacks with  unrealistic assumptions - e.g., acquiring privileged accounts

### Web
- Vulnerabilities in third-party tools
- Best practice concerns without PoC
- Clickjacking, open redirects (without impact)
- TLS config, SPF/DMARC/DNS misconfigs
- Lack of HTTP headers, verbose errors, self-XSS
- DoS/DDoS, social engineering, or phishing
- Issues only affecting outdated browsers
- Vulnerabilities requiring unlikely user actions

See the full list on our [Kaia Web HackenProof Program page](https://hackenproof.com/programs/kaia-web).

---

##  Rules & Guidelines

To participate, you must follow these rules:

- Test only in scope — no attacks on infrastructure or third-party systems
- Do not spam forms or create high-traffic scans
- Do not attempt DoS, phishing, or social engineering
- Do not access or modify other users’ data
- Do not disclose vulnerabilities publicly without our permission

All tests should be confined to your own accounts or test environments.

---

## Eligibility for Bounties

To qualify for a reward:
- Be the **first** to report the issue
- Submit only through HackenProof
- Include clear **steps to reproduce** and a **working PoC**
- Do not be a current/former employee or contractor

> AI-generated reports without a working PoC are **not eligible**.

---

## Coordinated Disclosure

- **Do not share vulnerabilities publicly**, even after they are resolved
- All communication must go through HackenProof
- Public disclosure will disqualify the report

---

## Resources

-  [Kaia Bounty Program on HackenProof](https://hackenproof.com/companies/kaia)  
-  [Kaia Documentation](https://docs.kaia.io)

---

Thanks for helping us improve Kaia’s security! We appreciate every responsible disclosure.
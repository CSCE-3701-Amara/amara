# Backend Modules

The Go backend is a modular monolith. Each directory is a business module with its own domain rules, application use cases, ports, adapters, and transport code as those features are implemented.

- `auth`: accounts and access control
- `cand`: candidate-facing experiences
- `job`: postings and application admission
- `invw`: interview workflow, assessments, and offers
- `ana`: analytics and evaluation

`platform` contains shared technical infrastructure. `integration` contains external provider adapters. Neither directory owns business state.
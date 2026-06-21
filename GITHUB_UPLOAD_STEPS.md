# GitHub Upload Steps

## Recommended Repo Name

```text
go-web-app-devops-portfolio
```

## Recommended Repo Description

```text
Go web application DevOps portfolio with Docker, GitHub Actions, Trivy, CodeQL, SBOM generation, Helm Kubernetes deployment, health checks, metrics, and optional AI release summaries.
```

## Create Repository

When creating the GitHub repository, keep these unchecked because this project already includes them:

```text
Add README file      unchecked
Add .gitignore       unchecked
Choose license       unchecked
```

## Upload Using Git Commands

```bash
git init
git add .
git commit -m "Initial commit - Go web app DevOps portfolio"
git branch -M main
git remote add origin https://github.com/YOUR_USERNAME/go-web-app-devops-portfolio.git
git push -u origin main
```

## Upload Using GitHub Website

If GitHub blocks large drag-and-drop uploads, upload in batches:

1. Root files: `README.md`, `Dockerfile`, `docker-compose.yaml`, `go.mod`, `main.go`, `main_test.go`, `.gitignore`, `.dockerignore`
2. `.github` folder
3. `static` folder
4. `helm` folder
5. `docs` folder
6. `scripts` folder

Then click **Commit changes** after each batch.

#!/usr/bin/env Rscript --vanilla

options(repos = c(REPO_NAME = "https://packagemanager.rstudio.com/all/__linux__/focal/latest"))

install.packages(c("jsonlite", "RPostgres"))

if (!requireNamespace("BiocManager", quietly = TRUE)) {
  install.packages("BiocManager")
}

databases <- BiocManager::available("^org\\..*\\..*\\.db$", include_installed = TRUE)
BiocManager::install(databases, update = FALSE, ask = FALSE)

BiocManager::install(c("gage", "pathview"), update = FALSE, ask = FALSE)

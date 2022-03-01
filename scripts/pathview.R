#!/usr/bin/env Rscript --vanilla

# .libPaths("/opt/R/library")
# .libPaths("/data/R")

require(jsonlite, quietly = TRUE)
require(pathview, quietly = TRUE)

### defaults
max_pathways <- 20

args <- commandArgs(TRUE)
json <- read_json(args[1])

con <- DBI::dbConnect(RPostgres::Postgres(),
        host = 'postgres',
        user = Sys.getenv("POSTGRES_USER"),
        password = Sys.getenv("POSTGRES_PASSWORD"),
        dbname = Sys.getenv("POSTGRES_DB"))

read_file <- function(json) {
    delim <- switch(tools::file_ext(json[['file']]),
        "csv" = ",",
        "tsv" = "\t",
        "txt" = "\t"
    )
    if (is.null(delim)) {
        stop("unsupported file type")
        return()
    }

    data <- read.delim(json[['file']], sep=delim, row.names=NULL)

    if (ncol(data) == 1) {
        d <- data
    } else if (ncol(data) > 1) {
        d <- data[,-1]
        rownames(d) <- data[,1]

        if (!is.null(json[["paired"]]) && json[["paired"]]) {
            if (length(json[['reference_col']]) != length(json[['sample_col']])) {
                stop("number of reference and sample columns not equal and paired mode enabled")
            }
        } else if (length(json[['reference_col']]) == 1) {

        } else {

        }
    } else {
        stop("empty file?")
    }

    return(d)
}

gene_data     <- read_file(json[['gene']])
compound_data <- read_file(json[['compound']])

if (json[['match']] && ncol(gene_data) != ncol(compound_data)) {
    stop("number of gene and compound columns not equal and matched mode enabled")
}

### pathway selection
if (is.null(json[['pathway']]) || json[['pathway']][['auto']]) {

    if (is.null(json[['pathway']][['ids']])) {
        stop("no pathways auto selected")
    }
} else if (is.null(json[['pathway']][['ids']])) {
    stop("no pathways provided")
}

json[['pathway']][['ids']] <- unlist(json[['pathway']][['ids']])
json[['pathway']][['ids']] <- unique(json[['pathway']][['ids']])
json[['pathway']][['ids']] <- json[['pathway']][['ids']][!is.na(json[['pathway']][['ids']])]
res <- DBI::dbSendQuery(con, "select path_id from kegg.organism_pathway where org_id = (select id from kegg.organism where code = $1 limit 1);")
DBI::dbBind(res, json[['species']])
organism_pathways <- unlist(DBI::dbFetch(res))
invalid_pathways <- setdiff(json[['pathway']][['ids']], organism_pathways)
json[['pathway']][['ids']] <- intersect(json[['pathway']][['ids']], organism_pathways)
DBI::dbClearResult(res)

if (length(json[['pathway']][['ids']]) > max_pathways) {
    skipped_pathways <- json[['pathway']][['ids']][max_pathways + 1:length(json[['pathway']][['ids']])]
    json[['pathway']][['ids']] <- json[['pathway']][['ids']][1:max_pathways]
}

if (!is.null(invalid_pathways) && length(invalid_pathways) > 0) {
    write.table(invalid_pathways, file="invalid-pathways.txt", sep="\t", col.names=NA, quote = FALSE)
}

if (exists("skipped_pathways") && !is.null(skipped_pathways) && length(skipped_pathways) > 0) {
    write.table(skipped_pathways, file="skipped-pathways.txt", sep="\t", col.names=NA, quote = FALSE)
}

if (length(json[['pathway']][['ids']]) == 0) {
    stop("no pathways selected after filtering")
}

### pathview arguments
pathview_args <- list(
    gene.data = gene_data,
    gene.idtype = json[['gene']][['id_type']],
    cpd.data = compound_data,
    cpd.idtype = json[['compound']][['id_type']],
    species = json[['species']],
    kegg.dir = Sys.getenv("KEGG_DIR"),
    match.data = json[['match']],
    out.suffix = json[['output']][['suffix']],
    sign.pos = "bottomright" # only applies when kegg.native = FALSE
)

if (!is.null(json[['output']][['format']]) && json[['output']][['format']] == "pdf") {
    pathview_args <- c(pathview_args, kegg.native = FALSE)
} else {
    pathview_args <- c(pathview_args, kegg.native = TRUE)
}

# same.layer = XXX,
# keys.align = XXX,
# split.group = XXX,
# expand.node = XXX,
# multi.state = XXX,
# node.sum = XXX,
# key.pos = XXX,
# cpd.label.offset = XXX,
# limit = list(gene = XXX, cpd = XXX),
# bins = list(gene = XXX, cpd = XXX),
# low = list(gene = XXX, cpd = XXX),
# mid = list(gene = XXX, cpd = XXX),
# high = list (gene = XXX, cpd = XXX),
# discrete = list(gene = XXX, cpd = XXX)

eat <- sapply(json[['pathway']][['ids']], function(pathway) {
    tryCatch(
        {
            output <- do.call(pathview, c(pathview_args, pathway.id = sprintf("%05d", pathway)))

            # if (json[['output']][['format']] == "png") {
            #     pv.labels(pv.out=output, pv.data.type=datatypes, pid=pathway)
            # }

            return()
        },
        error=function(cond) {
            message(cond)
            return()
        },
        warning=function(cond) {
            message(cond)
            return()
        }
    )

    return()
})

DBI::dbDisconnect(con)

warnings()

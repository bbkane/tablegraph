# Print lines/added deleted in a git repo over time

# https://git-scm.com/docs/git-log#Documentation/git-log.txt-emHem
# %H - commit hash
# %x09 - tab character
# %aI - author date
#
# RS = "" means it'll use a blank line as the field separator
git log \
    --format=format:"%aI" \
    --reverse \
    --shortstat \
| awk '
    BEGIN { RS = ""; FS = "\n"; OFS="\t"; print "date", "type", "lines" }
    {
        insertions = match($2, /[[:digit:]]+ insertion/)
        if (insertions != 0)
        {
            insertions = substr($2, RSTART, RLENGTH - 10)
            print $1, "insertion", insertions
        }

        deletions = match($2, /[[:digit:]]+ deletion/)
        if (deletions != 0)
        {
            deletions = -substr($2, RSTART, RLENGTH - 9)
            print $1, "deletion", deletions
        }
    }
' | tablegraph graph \
    --fieldsep $'\t' \
    --fieldnames firstline \
    --graph-title "Git History" \
    --type stacked-bar \
    --x-scale-type utc \
    --x-time-unit utcyearmonthdate \
    --x-type temporal \
    --y-type quantitative \
    --mark-size 5 \
| open_tmp_html.py

# TODO: GitHub stats with starghaze

# osquery

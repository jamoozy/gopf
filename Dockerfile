FROM scratch

COPY ./static /gopf
COPY index.tmpl.html /
COPY *.js /
COPY *.css /

CMD ["/gopf"]

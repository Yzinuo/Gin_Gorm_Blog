FROM mysql:8.0

ENV WORK_PATH  /usr/local/work
ENV AUTO_RUN_DIR /docker-entrypoint-initdb.d
ENV RUN_SHELL run.sh

COPY ./mysql.cnf /etc/mysql/conf.d
COPY ./gvb.sql ${WORK_PATH}/
COPY ./${RUN_SHELL} ${AUTO_RUN_DIR}

RUN chmod a+x ${AUTO_RUN_DIR}/${RUN_SHELL}
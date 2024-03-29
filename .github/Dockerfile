FROM myoung34/github-runner:latest

RUN set -euxo pipefail \
    && curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add - \
    && echo "deb https://dl.yarnpkg.com/debian/ stable main" > /etc/apt/sources.list.d/yarn.list \
    && apt update \
    && apt install --no-install-recommends -y \
        openssl \
        lsof \
        zip \
        make \
        curl \
        git \
        nodejs \
        yarn \
        postgresql-client \
    && yarn global add n \
    && yarn config set networkConcurrency 100 \
    && n v21 \
    && curl -sSL https://git.io/g-install | HOME=/home/runner sh -s -- -y bash \
    && ln -s /home/runner/go/bin/* /usr/local/bin/ \
    && rm -rf /var/lib/apt/lists/* \
    && rm -rf /tmp/*

# error:25066067:DSO support routines:dlfcn_load:could not load the shared library
# https://github.com/nodejs/node/issues/43132
RUN sed -ie '/providers =/ s/^#*/#/' /etc/ssl/openssl.cnf
FROM rust:latest as build

WORKDIR /usr/src/stelo-bets
COPY . .

RUN cargo build --release

FROM gcr.io/distroless/cc-debian11

COPY --from=build usr/src/stelo-bets/target/release/stelo-bets /usr/local/bin/

WORKDIR /usr/local/bin
CMD ["stelo-bets"]
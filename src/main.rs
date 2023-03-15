use actix_web::{post, web, App, HttpResponse, HttpServer, Responder};
use reqwest::Client;
use serde::{self, Deserialize};
use validator::Validate;

#[actix_web::main]
async fn main() -> std::io::Result<()> {
	let api_key = std::env::var("API_KEY").expect("Missing API_KEY env var");
	let wallet_key = std::env::var("WALLET_KEY").expect("Missing WALLET_KEY env var");

	let keys = web::Data::new((wallet_key, api_key));

	let http_client = web::Data::new(reqwest::Client::new());

	HttpServer::new(move || {
		App::new()
			.app_data(keys.clone())
			.app_data(http_client.clone())
			.service(post_transaction)
	})
	.bind(("0.0.0.0", 8080))?
	.shutdown_timeout(10)
	.run()
	.await
}

#[derive(Deserialize, Debug)]
struct TxQuery {
	key: String,
}

#[derive(Deserialize, Debug)]
struct TxBody {
	wallet_id: i64,
	memo: Option<String>,
	assets: SteloAsset,
}

#[derive(Deserialize, Debug, Validate)]
#[serde(deny_unknown_fields)]
struct SteloAsset {
	#[validate(range(min = 1000, max = 1000000))]
	stelo: u64,
}

#[post("/tx")]
async fn post_transaction(
	http_client: web::Data<Client>,
	keys: web::Data<(String, String)>,
	query: web::Query<TxQuery>,
	body: web::Json<TxBody>,
) -> impl Responder {
	if query.key != keys.1.to_string() {
		return HttpResponse::Unauthorized();
	};

	// Check if stelo is in valid range
	if let Err(_) = body.assets.validate() {
		return HttpResponse::BadRequest();
	}

	// Check if their memo is valid
	let chance_input = match &body.memo {
		Some(memo_string) => match memo_string.parse::<u8>() {
			Ok(memo_uint) => match memo_uint {
				2..=100 => memo_uint,
				_ => return HttpResponse::BadRequest(),
			},
			Err(_) => match memo_string.as_str() {
				"" => 2,
				_ => return HttpResponse::BadRequest(),
			},
		},
		None => 2,
	};

	// Calculate their chance, including 5% house edge
	let edged_chance = (1f32 / f32::from(chance_input)) - 0.05f32 / f32::from(chance_input);

	let roll: f32 = rand::random();

	// Less than, because the range is 0..1 not 0..=1
	if roll < edged_chance {
		// Send their winnings on another worker so webhook
		// gets the 200 response
		actix_rt::spawn(async move {
			let result = http_client
				.post("https://api.stelo.finance/wallet/transactions")
				.header("Authorization", format!("wallet {}", keys.0.to_string()))
				.json(&serde_json::json!({
					"recipient": body.wallet_id.to_string(),
					"type": 2,
					"memo": format!("You won {}x, congrats!", chance_input),
					"assets": {
						"stelo": body.assets.stelo * u64::from(chance_input)
					}
				}))
				.send()
				.await;

			if let Err(err) = result {
				println!("Failed to send winnings: {}", err);
			}
		});
	}

	HttpResponse::Ok()
}

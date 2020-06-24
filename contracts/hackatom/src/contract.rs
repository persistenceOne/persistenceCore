use schemars::JsonSchema;
use serde::{Deserialize, Serialize};
use serde_json::Result;
use cosmwasm_std::{
    from_slice, generic_err, log, not_found, to_binary, to_vec, unauthorized, AllBalanceResponse,
    Api, BankMsg, Binary, CanonicalAddr, Env, Extern, HandleResponse, HumanAddr, InitResponse,
    MigrateResponse, Querier, QueryResponse, StdResult, Storage, Coin, CosmosMsg,
};

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
pub struct InitMsg {
    pub verifier: HumanAddr,
    pub beneficiary: HumanAddr,
}

/// MigrateMsg allows a priviledged contract administrator to run
/// a migration on the contract. In this (demo) case it is just migrating
/// from one hackatom code to the same code, but taking advantage of the
/// migration step to set a new validator.
///
/// Note that the contract doesn't enforce permissions here, this is done
/// by blockchain logic (in the future by blockchain governance)
#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
pub struct MigrateMsg {
    pub verifier: HumanAddr,
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
pub struct State {
    pub verifier: CanonicalAddr,
    pub beneficiary: CanonicalAddr,
    pub funder: CanonicalAddr,
}

// failure modes to help test wasmd, based on this comment
// https://github.com/cosmwasm/wasmd/issues/8#issuecomment-576146751
#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
#[serde(rename_all = "snake_case")]
pub enum HandleMsg {
    AssetMint {
        properties: String
    },

}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
#[serde(rename_all = "snake_case")]
pub enum QueryMsg {
    // returns a human-readable representation of the verifier
    // use to ensure query path works in integration tests
    Verifier {},
    // This returns cosmwasm_std::AllBalanceResponse to demo use of the querier
    OtherBalance { address: HumanAddr },
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
pub struct VerifierResponse {
    pub verifier: HumanAddr,
}

pub static CONFIG_KEY: &[u8] = b"config";

pub fn init<S: Storage, A: Api, Q: Querier>(
    deps: &mut Extern<S, A, Q>,
    env: Env,
    msg: InitMsg,
) -> StdResult<InitResponse> {
    deps.storage.set(
        CONFIG_KEY,
        &to_vec(&State {
            verifier: deps.api.canonical_address(&msg.verifier)?,
            beneficiary: deps.api.canonical_address(&msg.beneficiary)?,
            funder: env.message.sender,
        })?,
    );

    // This adds some unrelated data and log for testing purposes
    Ok(InitResponse {
        data: Some(vec![0xF0, 0x0B, 0xAA].into()),
        log: vec![log("Let the", "hacking begin")],
        ..InitResponse::default()
    })
}

pub fn migrate<S: Storage, A: Api, Q: Querier>(
    deps: &mut Extern<S, A, Q>,
    _env: Env,
    msg: MigrateMsg,
) -> StdResult<MigrateResponse> {
    let data = deps
        .storage
        .get(CONFIG_KEY)
        .ok_or_else(|| not_found("State"))?;
    let mut config: State = from_slice(&data)?;
    config.verifier = deps.api.canonical_address(&msg.verifier)?;
    deps.storage.set(CONFIG_KEY, &to_vec(&config)?);
    Ok(MigrateResponse::default())
}

pub fn handle<S: Storage, A: Api, Q: Querier>(
    deps: &mut Extern<S, A, Q>,
    env: Env,
    msg: HandleMsg,
) -> StdResult<HandleResponse<PersistenceSDK>> {
    match msg {
        HandleMsg::AssetMint { properties} => do_asset_mint(deps, env, properties),
    }
}

fn do_release<S: Storage, A: Api, Q: Querier>(
    deps: &mut Extern<S, A, Q>,
    env: Env,
) -> StdResult<HandleResponse> {
    let data = deps
        .storage
        .get(CONFIG_KEY)
        .ok_or_else(|| not_found("State"))?;
    let state: State = from_slice(&data)?;

    if env.message.sender == state.verifier {
        let to_addr = deps.api.human_address(&state.beneficiary)?;
        let from_addr = deps.api.human_address(&env.contract.address)?;
        let balance = deps.querier.query_all_balances(&from_addr)?;

        let res = HandleResponse {
            log: vec![log("action", "release"), log("destination", &to_addr)],
            messages: vec![BankMsg::Send {
                from_address: from_addr,
                to_address: to_addr,
                amount: balance,
            }
                .into()],
            data: None,
        };
        Ok(res)
    } else {
        Err(unauthorized())
    }
}

/// TerraMsg is an override of CosmosMsg::Custom to add support for Terra's custom message types
#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
#[serde(rename_all = "snake_case")]
pub struct AssetMintRaw {
    from: HumanAddr,
    chainID: String,
    maintainersID: String,
    classificationID: String,
    properties: String,
    lock: i64,
    burn: i64,
}

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
#[serde(rename_all = "snake_case")]
pub struct PersistenceSDK {
    msgtype: String,
    raw: AssetMintRaw,
}



// {"mint":{"msgtype":"assetFactory/mint","raw":""}}
// this is a helper to be able to return these as CosmosMsg easier
impl Into<CosmosMsg<PersistenceSDK>> for PersistenceSDK {
    fn into(self) -> CosmosMsg<PersistenceSDK> {
        CosmosMsg::Custom(self)
    }
}


fn do_asset_mint<S: Storage, A: Api, Q: Querier>(
    deps: &mut Extern<S, A, Q>,
    env: Env,
    properties: String,
) -> StdResult<HandleResponse<PersistenceSDK>> {
    let data = deps
        .storage
        .get(CONFIG_KEY)
        .ok_or_else(|| not_found("State"))?;
    let state: State = from_slice(&data)?;

    if env.message.sender == state.verifier {
        let from_addr = deps.api.human_address(&env.contract.address)?;

        let mintMsg = AssetMintRaw {
            from: deps.api.human_address(&env.message.sender)?,
            chainID: "test".to_owned(),
            maintainersID: "myid".to_owned(),
            classificationID: "classid".to_owned(),
            properties: properties,
            lock: -1,
            burn: -1,
        };

        let res = HandleResponse {
            log: vec![log("action", "asset_mint"), log("destination", &from_addr)],
            messages: vec![PersistenceSDK {
                msgtype: "assetFactory/mint".to_string(),
                raw: mintMsg,
            }
                .into()],
            data: None,
        };
        Ok(res)
    } else {
        Err(unauthorized())
    }
}



fn do_allocate_large_memory() -> StdResult<HandleResponse> {
    // We create memory pages explicitely since Rust's default allocator seems to be clever enough
    // to not grow memory for unused capacity like `Vec::<u8>::with_capacity(100 * 1024 * 1024)`.
    // Even with std::alloc::alloc the memory did now grow beyond 1.5 MiB.

    #[cfg(target_arch = "wasm32")]
    {
        use core::arch::wasm32;
        let pages = 1_600; // 100 MiB
        let ptr = wasm32::memory_grow(0, pages);
        if ptr == usize::max_value() {
            return Err(generic_err("Error in memory.grow instruction"));
        }
        Ok(HandleResponse::default())
    }

    #[cfg(not(target_arch = "wasm32"))]
    Err(generic_err("Unsupported architecture"))
}

fn do_panic() -> StdResult<HandleResponse> {
    panic!("This page intentionally faulted");
}

pub fn query<S: Storage, A: Api, Q: Querier>(
    deps: &Extern<S, A, Q>,
    msg: QueryMsg,
) -> StdResult<QueryResponse> {
    match msg {
        QueryMsg::Verifier {} => query_verifier(deps),
        QueryMsg::OtherBalance { address } => query_other_balance(deps, address),
    }
}

fn query_verifier<S: Storage, A: Api, Q: Querier>(
    deps: &Extern<S, A, Q>,
) -> StdResult<QueryResponse> {
    let data = deps
        .storage
        .get(CONFIG_KEY)
        .ok_or_else(|| not_found("State"))?;
    let state: State = from_slice(&data)?;
    let addr = deps.api.human_address(&state.verifier)?;
    Ok(Binary(to_vec(&VerifierResponse { verifier: addr })?))
}

fn query_other_balance<S: Storage, A: Api, Q: Querier>(
    deps: &Extern<S, A, Q>,
    address: HumanAddr,
) -> StdResult<QueryResponse> {
    let amount = deps.querier.query_all_balances(address)?;
    to_binary(&AllBalanceResponse { amount })
}

use schemars::JsonSchema;
use serde::{Deserialize, Serialize};

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
#[serde(rename_all = "snake_case")]
pub enum QueryMsg {
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

pub fn query<S: Storage, A: Api, Q: Querier>(
    deps: &Extern<S, A, Q>,
    msg: QueryMsg,
) -> StdResult<QueryResponse> {
    match msg {

    }
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
        messages: vec![],
        log: vec![log("Let the", "hacking begin")],
        ..InitResponse::default()
    })
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



// {"mint":{"msgtype":"assets/mint","raw":""}}
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

        // can add all the parameters as input params
        let mintMsg = AssetMintRaw {
            from: deps.api.human_address(&env.message.sender)?,
            chainID: "".to_owned(),
            maintainersID: "".to_owned(),
            classificationID: "".to_owned(),
            properties: properties,
            lock: -1,
            burn: -1,
        };

        let res = HandleResponse {
            log: vec![log("action", "asset_mint"), log("destination", &from_addr)],
            messages: vec![PersistenceSDK {
                msgtype: "assets/mint".to_string(),
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

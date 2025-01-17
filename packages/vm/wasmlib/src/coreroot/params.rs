// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

// (Re-)generated by schema tool
// >>>> DO NOT CHANGE THIS FILE! <<<<
// Change the json schema instead

#![allow(dead_code)]
#![allow(unused_imports)]

use crate::*;
use crate::coreroot::*;
use crate::host::*;

#[derive(Clone, Copy)]
pub struct ImmutableDeployContractParams {
    pub(crate) id: i32,
}

impl ImmutableDeployContractParams {
    pub fn description(&self) -> ScImmutableString {
        ScImmutableString::new(self.id, PARAM_DESCRIPTION.get_key_id())
    }

    pub fn name(&self) -> ScImmutableString {
        ScImmutableString::new(self.id, PARAM_NAME.get_key_id())
    }

    pub fn program_hash(&self) -> ScImmutableHash {
        ScImmutableHash::new(self.id, PARAM_PROGRAM_HASH.get_key_id())
    }
}

#[derive(Clone, Copy)]
pub struct MutableDeployContractParams {
    pub(crate) id: i32,
}

impl MutableDeployContractParams {
    pub fn description(&self) -> ScMutableString {
        ScMutableString::new(self.id, PARAM_DESCRIPTION.get_key_id())
    }

    pub fn name(&self) -> ScMutableString {
        ScMutableString::new(self.id, PARAM_NAME.get_key_id())
    }

    pub fn program_hash(&self) -> ScMutableHash {
        ScMutableHash::new(self.id, PARAM_PROGRAM_HASH.get_key_id())
    }
}

#[derive(Clone, Copy)]
pub struct ImmutableGrantDeployPermissionParams {
    pub(crate) id: i32,
}

impl ImmutableGrantDeployPermissionParams {
    pub fn deployer(&self) -> ScImmutableAgentID {
        ScImmutableAgentID::new(self.id, PARAM_DEPLOYER.get_key_id())
    }
}

#[derive(Clone, Copy)]
pub struct MutableGrantDeployPermissionParams {
    pub(crate) id: i32,
}

impl MutableGrantDeployPermissionParams {
    pub fn deployer(&self) -> ScMutableAgentID {
        ScMutableAgentID::new(self.id, PARAM_DEPLOYER.get_key_id())
    }
}

#[derive(Clone, Copy)]
pub struct ImmutableRevokeDeployPermissionParams {
    pub(crate) id: i32,
}

impl ImmutableRevokeDeployPermissionParams {
    pub fn deployer(&self) -> ScImmutableAgentID {
        ScImmutableAgentID::new(self.id, PARAM_DEPLOYER.get_key_id())
    }
}

#[derive(Clone, Copy)]
pub struct MutableRevokeDeployPermissionParams {
    pub(crate) id: i32,
}

impl MutableRevokeDeployPermissionParams {
    pub fn deployer(&self) -> ScMutableAgentID {
        ScMutableAgentID::new(self.id, PARAM_DEPLOYER.get_key_id())
    }
}

#[derive(Clone, Copy)]
pub struct ImmutableFindContractParams {
    pub(crate) id: i32,
}

impl ImmutableFindContractParams {
    pub fn hname(&self) -> ScImmutableHname {
        ScImmutableHname::new(self.id, PARAM_HNAME.get_key_id())
    }
}

#[derive(Clone, Copy)]
pub struct MutableFindContractParams {
    pub(crate) id: i32,
}

impl MutableFindContractParams {
    pub fn hname(&self) -> ScMutableHname {
        ScMutableHname::new(self.id, PARAM_HNAME.get_key_id())
    }
}

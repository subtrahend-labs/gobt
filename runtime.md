- [x] get_delegates 
- [x] get-delegated
- [ ] get_neurons_lite 
- [ ] get_neuron_lite
- [x] get_neurons
- [x] get_neuron
- [x] get_subnet_info
- [x] get_subnets_info
- [x] get_subnet_info_v2
- [x] get_subnets_info_v2
- [x] get_subnet_hyperparams
- [x] get_all_dynamic_info
- [x] get_dynamic_info
- [x] get_all_metagraphs
- [x] get_metagraph
- [x] get_subnet_state
- [ ] get_network_lock_cost
- [x] get_selective_metagraph

    #[method(name = "subnetInfo_getAllMetagraphs")]
    fn get_all_metagraphs(&self, at: Option<BlockHash>) -> RpcResult<Vec<u8>>;
    #[method(name = "subnetInfo_getMetagraph")]
    fn get_metagraph(&self, netuid: u16, at: Option<BlockHash>) -> RpcResult<Vec<u8>>;


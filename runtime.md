- [x] get_delegates 
- [x] get-delegated
- [ ] get_neurons_lite 
- [ ] get_neuron_lite
- [x] get_neurons
- [x] get_neuron
- [ ] get_subnet_info
- [ ] get_subnets_info
- [ ] get_subnet_info_v2
- [ ] get_subnets_info_v2
- [ ] get_subnet_hyperparams
- [ ] get_all_dynamic_info
- [ ] get_dynamic_info
- [ ] get_all_metagraphs
- [x] get_metagraph
- [ ] get_subnet_state
- [ ] get_network_lock_cost
- [ ] get_selective_metagraph

    #[method(name = "subnetInfo_getAllMetagraphs")]
    fn get_all_metagraphs(&self, at: Option<BlockHash>) -> RpcResult<Vec<u8>>;
    #[method(name = "subnetInfo_getMetagraph")]
    fn get_metagraph(&self, netuid: u16, at: Option<BlockHash>) -> RpcResult<Vec<u8>>;


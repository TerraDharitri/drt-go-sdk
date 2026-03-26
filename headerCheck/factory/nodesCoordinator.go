package factory

import (
	"github.com/TerraDharitri/drt-go-chain-core/core"
	"github.com/TerraDharitri/drt-go-chain-core/data/endProcess"
	crypto "github.com/TerraDharitri/drt-go-chain-crypto"
	"github.com/TerraDharitri/drt-go-chain/common"
	"github.com/TerraDharitri/drt-go-chain/config"
	"github.com/TerraDharitri/drt-go-chain/dataRetriever/dataPool"
	"github.com/TerraDharitri/drt-go-chain/sharding/nodesCoordinator"
	"github.com/TerraDharitri/drt-go-sdk/data"
	"github.com/TerraDharitri/drt-go-sdk/disabled"
)

const (
	defaultSelectionChances = uint32(1)
)

// CreateNodesCoordinator creates nodes coordinator which will be used for header verification
func CreateNodesCoordinator(
	coreComp *coreComponents,
	networkConfig *data.NetworkConfig,
	enableEpochsConfig *data.EnableEpochsConfig,
	publicKey crypto.PublicKey,
	genesisNodesConfig *data.GenesisNodes,
	chainParams nodesCoordinator.ChainParametersHandler,
) (nodesCoordinator.EpochsConfigUpdateHandler, error) {
	eligibleValidators, err := generateGenesisNodes(coreComp.PubKeyConverter, genesisNodesConfig.Eligible)
	if err != nil {
		return nil, err
	}

	waitingValidators, err := generateGenesisNodes(coreComp.PubKeyConverter, genesisNodesConfig.Waiting)
	if err != nil {
		return nil, err
	}

	argsNodesShuffler := createArgsNodesShuffler(enableEpochsConfig, coreComp.EnableEpochsHandler)
	nodeShuffler, err := nodesCoordinator.NewHashValidatorsShuffler(argsNodesShuffler)
	if err != nil {
		return nil, err
	}

	publicKeyBytes, err := publicKey.ToByteArray()
	if err != nil {
		return nil, err
	}

	initialEpoch := uint32(0)
	arguments := nodesCoordinator.ArgNodesCoordinator{
		
		ChainParametersHandler:          chainParams,
		GenesisNodesSetupHandler:        &disabled.GenesisNodesSetupHandler{},
		NodesCoordinatorRegistryFactory: &disabled.NodesCoordinatorRegistryFactory{},
		Marshalizer:                     coreComp.Marshaller,
		Hasher:                          coreComp.Hasher,
		Shuffler:                        nodeShuffler,
		EpochStartNotifier:              &disabled.EpochStartNotifier{},
		BootStorer:                      &disabled.Storer{},
		ShardIDAsObserver:               0,
		NbShards:                        networkConfig.NumShardsWithoutMeta,
		EligibleNodes:                   eligibleValidators,
		WaitingNodes:                    waitingValidators,
		SelfPublicKey:                   publicKeyBytes,
		Epoch:                           initialEpoch,
		StartEpoch:                      0,
		ConsensusGroupCache:             &disabled.Cache{},
		ShuffledOutHandler:              &disabled.ShuffledOutHandler{},
		ChanStopNode:                    make(chan endProcess.ArgEndProcess),
		NodeTypeProvider:                &disabled.NodeTypeProvider{},
		IsFullArchive:                   false,
		EnableEpochsHandler:             coreComp.EnableEpochsHandler,
		ValidatorInfoCacher:             dataPool.NewCurrentEpochValidatorInfoPool(),
	}

	baseNodesCoordinator, err := nodesCoordinator.NewIndexHashedNodesCoordinator(arguments)
	if err != nil {
		return nil, err
	}

	nd, err := nodesCoordinator.NewIndexHashedNodesCoordinatorWithRater(baseNodesCoordinator, coreComp.Rater)
	if err != nil {
		return nil, err
	}

	return nd, nil
}

func generateGenesisNodes(converter core.PubkeyConverter, nodesConfig map[uint32][]string) (map[uint32][]nodesCoordinator.Validator, error) {
	validatorsMap := make(map[uint32][]nodesCoordinator.Validator)

	for shardID, nodesPubKeys := range nodesConfig {
		validators := make([]nodesCoordinator.Validator, 0, len(nodesPubKeys))
		for i, pubKey := range nodesPubKeys {
			pubKeyBytes, err := converter.Decode(pubKey)
			if err != nil {
				return nil, err
			}

			validator, err := nodesCoordinator.NewValidator(pubKeyBytes, defaultSelectionChances, uint32(i))
			if err != nil {
				return nil, err
			}

			validators = append(validators, validator)
		}
		validatorsMap[shardID] = validators
	}

	return validatorsMap, nil
}

func createArgsNodesShuffler(
	eec *data.EnableEpochsConfig,
	enableEpochsHandler common.EnableEpochsHandler,
) *nodesCoordinator.NodesShufflerArgs {
	maxNodesChangeConfigs := make([]config.MaxNodesChangeConfig, 0)
	for _, conf := range eec.EnableEpochs.MaxNodesChangeEnableEpoch {
		maxNodesChangeConfig := config.MaxNodesChangeConfig{
			EpochEnable:            conf.EpochEnable,
			MaxNumNodes:            conf.MaxNumNodes,
			NodesToShufflePerShard: conf.NodesToShufflePerShard,
		}

		maxNodesChangeConfigs = append(maxNodesChangeConfigs, maxNodesChangeConfig)
	}

	argsNodesShuffler := &nodesCoordinator.NodesShufflerArgs{

		ShuffleBetweenShards: true,
		MaxNodesEnableConfig: maxNodesChangeConfigs,
		EnableEpochsHandler:  enableEpochsHandler,
		EnableEpochs:         eec.EnableEpochs,
	}

	return argsNodesShuffler
}

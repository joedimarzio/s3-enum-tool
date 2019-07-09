

package main

import (
    "github.com/aws/aws-sdk-go/aws"
    log "github.com/sirupsen/logrus"

    "slurp/scanner/external"
    "slurp/scanner/cmd"
    "slurp/scanner/intern"
)

// Global config
var cfg cmd.Config

func main() {
    cfg = cmd.Init("slurp", "Public buckets finder", "Public buckets finder")

    switch cfg.State {
    case "DOMAIN":
        var config aws.Config
        config.Region = &cfg.Region
        external.Init(&cfg)

        log.Info("Building permutations....")
        bucketNames := external.PermutateDomainRunner(&cfg) //go

        log.Info("Processing permutations....")
        external.CheckDomainPermutations(&cfg, config, bucketNames)

        // Print stats info
        log.Printf("%+v", cfg.Stats)
    case "INTERNAL":
        var config aws.Config
        config.Region = &cfg.Region

        log.Info("Determining public buckets....")
        buckets, err3 := intern.GetPublicBuckets(config)
        if err3 != nil {
            log.Error(err3)
        }

        for bucket := range buckets.ACL {
            log.Infof("S3 public bucket (ACL): %s", buckets.ACL[bucket])
        }

        for bucket := range buckets.Policy {
            log.Infof("S3 public bucket (Policy): %s", buckets.Policy[bucket])
        }

    default:
        log.Fatal("Check help")
    }
}

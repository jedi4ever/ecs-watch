# Description

**CAVEAT!** this is WIP for now - soon to be released

When you want to run microservices inside of an ECS cluster , the easiest/standard/recommended way is to register Services with ELBs.
You can register the ELBs in DNS and each service is instanly available. No need for a service discovery.

As the number of microservices grow , so increases the cost of the ELBs. Especially in Test/Staging environment this can get quite costly where multiple environment are running on the same cluster.

The alternative is to use a discovery service and a reverse proxy like (Nginx, HaProxy, Traefik, Kong, Fabio).
Most tools/examples I found used Consul. But then I thought, why duplicate the ECS state in consul? That way you can avoid running a consul cluster.

## Report
```
NAME:
   ecs-watch report - reports all containers and ports

USAGE:
   ecs-watch report [command options] [arguments...]

OPTIONS:
   --cluster value  (default: "default") [$ECSWATCH_CLUSTER]
   --region value   (default: "eu-west-1") [$ECSWATCH_REGION]
```

```
+-----------+----------+---------------+----------------+-----------+---------------------+-----------+-------------+---------+--------------------+
|   NAME    | HOSTPORT | CONTAINERPORT |    PUBLICIP    | PRIVATEIP |     INSTANCEID      |   IMAGE   | VIRTUALHOST | STATUS  |      CLUSTER       |
+-----------+----------+---------------+----------------+-----------+---------------------+-----------+-------------+---------+--------------------+
| jenkins-b |    32768 |            80 | 52.209.248.211 | 10.0.4.42 | i-08d46a56b977b8a62 | httpd:2.4 |             | RUNNING | staging-asg        |
| jenkins-a |    32769 |            80 | 52.209.248.211 | 10.0.4.42 | i-08d46a56b977b8a62 | httpd:2.4 |             | RUNNING | staging-asg        |
+-----------+----------+---------------+----------------+-----------+---------------------+-----------+-------------+---------+--------------------+
```

## Generate file using Template
The use of *elb-watch* to have a command similar to docker-gen/nginx-proxy where a template is populated and a signal is send to a container.
instead of reading the state from Consul, it reads the state from ECS.

```
NAME:
   ecs-watch generate - generates a file

USAGE:
   ecs-watch generate [command options] [arguments...]

OPTIONS:
   --cluster value           (default: "default") [$ECSWATCH_CLUSTER]
   --region value            (default: "eu-west-1") [$ECSWATCH_REGION]
   --template-file value      [$ECSWATCH_TEMPLATE_FILE]
   --output-file value        [$ECSWATCH_OUTPUT_FILE]
   --notify-container value   [$ECSWATCH_NOTIFY_CONTAINER]
   --docker-signal value     (default: "SIGHUP") [$ECSWATCH_DOCKER_SIGNAL]
   --docker-container value   [$ECSWATCH_DOCKER_CONTAINER]
   --docker-endpoint value   (default: "unix:///var/run/docker.sock") [$ECSWATCH_DOCKER_ENDPOINT]
   --watch value             (default: "false") [$ECSWATCH_WATCH]
```

## Route53 sidekick
The use of *elb-watch* to find the connection details of container (ip/port) and update a record in route53.
This is useful as a side-kick for services that are not http based (for example your redis server)

# Related projects
- <https://github.com/CpuID/ecs-discoverer>
- <https://github.com/majest/docker-consul-ecs>
- <https://github.com/adamdecaf/aws-ecs-nginx-proxy>
- <https://github.com/awslabs/service-discovery-ecs-consul>
- <https://github.com/awslabs/ecs-refarch-service-discovery/>
- <https://github.com/kyani-inc/ecs-discovery>
- <https://github.com/gliderlabs/registrator>
- <https://github.com/jwilder/docker-gen>

# Reading
- <https://aws.amazon.com/blogs/compute/service-discovery-via-consul-with-amazon-ecs/>
- <https://aws.amazon.com/blogs/compute/service-discovery-an-amazon-ecs-reference-architecture/>
- <http://www.slideshare.net/AmazonWebServices/microservices-and-amazon-ecs>
- <https://segment.com/blog/rebuilding-our-infrastructure/>
- <https://github.com/rhockenbury/trace>
- <https://github.com/jhspaybar/ecs_state>
- <http://engineering.skybettingandgaming.com/2016/05/05/aws-and-consul/>
- <https://sitano.github.io/2015/10/06/abt-consul-outage/>
- <https://www.nginx.com/blog/service-discovery-in-a-microservices-architecture/>

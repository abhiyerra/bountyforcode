if(window.location.href.match("localhost")) {
    BountyForCode = {
        server: "http://localhost:3000",
        ghWatch: "http://ghwatch.bountyforcode.com"
    }
} else if(window.location.href.match("beta")) {
    BountyForCode = {
        server: "http://betaapi.bountyforcode.com",
        ghWatch: "http://ghwatch.bountyforcode.com"
    }
} else {
    BountyForCode = {
        server: "http://api.bountyforcode.com",
        ghWatch: "http://ghwatch.bountyforcode.com"
    }
}

angular.module('bountyforcode', ['ngRoute'])
    .config(function($routeProvider) {
        $routeProvider
            .when('/', {
                controller:'WelcomeCtrl',
                templateUrl:'templates/issues/welcome.html'
            })
            .when('/issues/:issueId', {
                controller:'ShowIssueCtrl',
                templateUrl:'templates/issues/show.html'
            })
            .when('/issues/:issueId/contribute', {
                controller:'ContributeIssueCtrl',
                templateUrl:'templates/issues/contribute.html'
            })
            .when('/issues/:issueId/claim', {
                controller:'ClaimBountyCtrl',
                templateUrl:'templates/bounties/claim.html'
            })
            .otherwise({
                redirectTo:'/'
            });
    })

    .factory('BFCApiService', function($http) {
        return {
            getIssues: function() {
                return $http.get(BountyForCode.server + "/v1/issues")
                    .then(function(response) {
                        console.log("Issues: " + response.data);
                        return response.data;
                    });
            },

            getIssue: function(issueId) {
                return $http.get(BountyForCode.server + "/v1/issues/" + issueId)
                    .then(function(response) {
                        console.log("Issues: " + response.data);
                        return response.data;
                    });
            },

            postIssue: function(issueUrl) {
                return $http.post(
                    BountyForCode.server + "/v1/issues",
                    {"issue_url": issueUrl},
                    { headers: {'Content-Type': 'application/x-www-form-urlencoded'} }
                ).then(function(response) {
                    console.log("New Issue: " + response.data);
                    return response.data;
                });
            },

            newBounty: function(issueId) {
                return $http.get(BountyForCode.server + "/v1/issues/" + issueId + "/bounty")
                    .then(function(response) {
                        console.log("New Bounty: " + response.data);
                        return response.data;
                    });
            },

            getBounties: function(issueId) {
                return $http.get(BountyForCode.server + "/v1/issues/" + issueId + "/bounties")
                    .then(function(response) {
                        console.log("Get Bounties: " + response.data);
                        return response.data;
                    });

            },

            getUser: function() {
                return $http.get(BountyForCode.server + "/v1/user")
                    .then(function(response) {
                        console.log("User: " + response.data);
                        return response.data;
                    });

            },

            getUserForId: function(id) {
                return $http.get(BountyForCode.server + "/v1/users/" + id)
                    .then(function(response) {
                        console.log("User: " + response.data);
                        return response.data;
                    });

            },

            getRegisterUrl: function() {
                return BountyForCode.server + "/v1/register?redirect=" + escape(window.location.href);
            }
        }
    })

    .factory('GithubApiService', function($http) {
        return {
            getIssue: function(project, repo, id) {
                return $http.get(BountyForCode.ghWatch + "/v1/issues?issue_url=http://github.com/" + project + "/" + repo + "/issues/" + id)
                    .then(function(response) {
                        console.log("Github Issue: " + response.data);
                        return response.data;
                    });
            },

            getIssueUrl: function(url) {
                return $http.get(BountyForCode.ghWatch + "/v1/issues?issue_url=" + url)
                    .then(function(response) {
                        console.log("Github Issue: " + response.data);
                        return response.data;
                    });
            },

            getUser: function(username) {
                return $http.get(BountyForCode.ghWatch + "/v1/users?github_login=" + username)
                    .then(function(response) {
                        console.log("Github Username: " + response.data);
                        return response.data;
                    });
            }
        }
    })

    .controller('HeaderCtrl', function($scope, $location, BFCApiService) {
        $scope.showHeader = true;

        $scope.registerUrl = BFCApiService.getRegisterUrl();

        BFCApiService.getUser().then(function(data) {
            $scope.loginText = "Login with Github";
        });

    })

    .controller('WelcomeCtrl', function($scope) {
        $scope.showHeader = false;
    })

    .controller('NewBountyCtrl', function($scope, $location, $http, $sce, BFCApiService, GithubApiService) {
        $scope.issue = {};

        $scope.getIssue = function() {
            $scope.issue = {};
            $scope.desc = "";

            GithubApiService.getIssueUrl($scope.issueUrl).then(function(issue) {
                $scope.issue = issue;
                $scope.desc = $sce.trustAsHtml(markdown.toHTML(issue.body));

            });
        };

        $scope.create = function() {
            $scope.message = "";

            BFCApiService.postIssue(this.issueUrl).then(function(data) {
                $scope.message = data;
            });
        };
    })

    .controller('ShowIssueCtrl', function($scope, $routeParams, $sce, BFCApiService, GithubApiService) {
        $scope.issue = {};

        $scope.githubIssue = {};

        BFCApiService.getIssue($routeParams.issueId).then(function(issue) {
            $scope.issue = issue;

            GithubApiService.getIssue($scope.issue.project, $scope.issue.repo, $scope.issue.identifier).then(function(issue) {
                $scope.github = issue;
                if(issue.body != null) {
                    $scope.desc = $sce.trustAsHtml(markdown.toHTML(issue.body));
                }
            });
        });
    })

    .controller('IssueBountyCtrl', function($scope, $http, BFCApiService, GithubApiService) {

        BFCApiService.getUserForId($scope.bounty.user_id).then(function(user) {
            $scope.user = user;

            GithubApiService.getUser(user.github_username).then(function(github_user) {
                $scope.github = github_user;
            });

        });
    })

    .controller('IssueBountiesCtrl', function($scope, $routeParams, $http, BFCApiService, GithubApiService) {
        BFCApiService.getBounties($routeParams.issueId).then(function(bounties) {
            $scope.bounties = bounties;
        });
    })

    .controller('DiscoverBountyCtrl', function($scope, $http, BFCApiService, GithubApiService) {
        $scope.issues = [];

        BFCApiService.getIssues().then(function(issues) {
            $scope.issues = issues;
        });
    })

    .controller('DiscoverIssueCtrl', function($scope, $http, BFCApiService, GithubApiService) {
        var project = $scope.issue.project;
        var repo = $scope.issue.repo;
        var identifier = $scope.issue.identifier;
        var id = $scope.issue.id;

        GithubApiService.getIssue(project, repo, identifier).then(function(github) {
            $scope.github = github;
        });
    })

    .controller('ContributeIssueCtrl', function($scope, $routeParams, BFCApiService) {
        $scope.issue = {
            title: "Hello",
            githubUrl: "http://github.com/abhiyerra/feedbackjs/issues/1",
        };
        $scope.bounty = {};
        $scope.issueId = $routeParams.issueId;

        BFCApiService.newBounty($routeParams.issueId).then(function(bounty) {
            $scope.bounty = bounty;
        });
    });

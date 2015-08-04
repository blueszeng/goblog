var blogAdminControllers = angular.module('blogAdminControllers', []);

blogAdminControllers.service('UserService', [
	function() {
		
  var currentUser = []

  var getUser = function (){
      return currentUser[0];
  };

  var saveUser = function(newObj) {
      currentUser.push(newObj);
  };
  
  return {
    saveUser: saveUser,
    getUser: getUser
  };

}]);

blogAdminControllers.controller('AdminHeaderCtrl', ['$scope', '$http', '$timeout', '$state', 'UserService',
	function($scope, $http, $timeout, $state, UserService) {

	$http.get('/api/users').success(function(data) {
		$timeout(function() {
			$scope.user = data;
			UserService.saveUser(data);
		}, 0);
	});

	$scope.$on('scopeChanged', function() {
		$scope.currentState = $state.current.name
		console.log($scope.currentState)
	});
}]);

blogAdminControllers.controller('AdminHomeCtrl', ['$rootScope', '$scope', '$http', '$state', '$timeout', '$sce', 'UserService',
	function($rootScope, $scope, $http, $state, $timeout, $sce, UserService) {


	$timeout(function() {
		$scope.user = UserService.getUser();
		$scope.loginURL = $sce.trustAsResourceUrl($scope.user.loginURL);
	}, 600);

	$scope.$watch('user', function() {
		if (angular.isUndefined($scope.user)) {
			$timeout(function() {
				$scope.user = UserService.getUser();
				if($scope.user["displayName"] == "" && $scope.user["role"] == "SiteAdmin") {
					console.log("New Site Admin")
					$state.go('root.useredit')
				}
			}, 600);
		}
	});
   
	$rootScope.$broadcast('scopeChanged', "root.home")
}]);

blogAdminControllers.controller('UsersListCtrl', ['$rootScope', '$scope', '$http', '$state', '$timeout', 'UserService',
	function($rootScope, $scope, $http, $state, $timeout, UserService) {

	$scope.user = UserService.getUser()

	$scope.$watch('user', function() {
		if (angular.isUndefined($scope.user)) {
			$state.go('root.home')	   	
		} else if ($scope.user["role"] != "SiteAdmin") {
			$state.go('root.home')
		}
	});
   
	$http.get('/api/users/all').success(function(data) {
		$timeout(function() {
			$scope.users = data
		}, 0);
	});
	
	$rootScope.$broadcast('scopeChanged', "root.users")
}]);

blogAdminControllers.controller('UserEditCtrl', ['$rootScope', '$scope', '$http', '$state', '$timeout', 'UserService',
	function($rootScope, $scope, $http, $state, $timeout, UserService) {

	$scope.user = UserService.getUser()

	$scope.$watch('user', function() {
		if (angular.isUndefined($scope.user)) {
			console.log("No User Logged In")
			$state.go('root.home')	   	
		} else if ($scope.user["email"] == "") {
			console.log("No User Email Address")
			$state.go('root.home')
		} else {
			console.log($scope.user["role"])
		}
	});

	$rootScope.$broadcast('scopeChanged', "root.home")
	
	$scope.update = function(user) {
	 	$http.post('/api/users', user).success(function() {
	        $timeout(function() {
	        	$state.go('root.home', {}, {reload: true});
	        }, 100);
		})
	};

	$scope.cancelEdit = function() {
		$state.go('root.home')
	}
	
	$scope.add = function(newUser) {
	 	$http.post('/api/users', newUser).success(function() {
	        $timeout(function() {
	            	$state.go('^', {}, {reload: true});
	        }, 100);
		})
	};

    	
}]);
	
		
blogAdminControllers.controller('BlogsListCtrl', ['$rootScope', '$scope', '$http', '$state', '$timeout', 'UserService',
	function($rootScope, $scope, $http, $state, $timeout, UserService) {

	$scope.user = UserService.getUser()

	$scope.$watch('user', function() {
		if (angular.isUndefined($scope.user)) {
			$state.go('root.home')	   	
		} else if ($scope.user["role"] != "SiteAdmin") {
			$state.go('root.home')
		}
	});
	
	$scope.blogLoad = function () {
		$http.get('/api/blogs/all').success(function(data) {
			$timeout(function() {
				$scope.blogs = data
			}, 0);
		});
	};
	
	$scope.blogLoad()

	$rootScope.$broadcast('scopeChanged', "root.blogs")
}]);

blogAdminControllers.controller('BlogEditCtrl', ['$scope', '$http', '$stateParams', '$state', '$timeout', '$location', '$anchorScroll',
	function($scope, $http, $stateParams, $state, $timeout, $location, $anchorScroll) {
    	$scope.blogID = $stateParams.blogID;

    	$http.get('/api/blogs/'+$scope.blogID).success(function(data) {
    
    		$scope.blog = data
	    	$timeout(function() {
	    		//$location.hash('blogEdit');
	    		//$anchorScroll(); 
	    	}, 100);
    	});

    	$scope.update = function(blog) {
	     	$http.post('/api/blogs', blog).success(function() {
	            $timeout(function() {
	            	$state.go('^', {}, {reload: true});
	            }, 100);
	     	})
      	};
      	
      	$scope.deleteEmail = function (index) {
        	$scope.blog.blogEmails.splice(index, 1);
        	$scope.blog.blogAuthors.splice(index, 1);
    	}
    	
    	$scope.addEmail = function (index) {
    		$http.get('/api/users/' + $scope.newEmail).success(function(data) {
        		$scope.blog.blogEmails.push($scope.newEmail);
        		$scope.blog.blogAuthors.push(data);    			
    		})
    	}
    	
    	$scope.cancelEdit = function() {
    		$state.go('^')
    	}
    }]);
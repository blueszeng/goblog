var blogAdminControllers = angular.module('blogAdminControllers', []);

blogAdminControllers.service('UserService', [
	function() {
		
  var currentUser = []

  var getUser = function (){
	console.log("Retrieve User")
	console.log(currentUser[0])
	return currentUser[0];
  };

  var saveUser = function(newObj) {
	currentUser.splice(0)
	//console.log(currentUser)
	console.log("Store User")
	currentUser.push(newObj);
  };
  
  return {
    saveUser: saveUser,
    getUser: getUser
  };

}]);

blogAdminControllers.controller('AdminHeaderCtrl', ['$rootScope', '$scope', '$http', '$timeout', '$state', 'UserService',
	function($rootScope, $scope, $http, $timeout, $state, UserService) {

	$http.get('/api/users').success(function(data) {
		$timeout(function() {
			$scope.user = data;
		}, 0);
		UserService.saveUser(data);
		$rootScope.$broadcast('scopeChanged', "root.home")
	});
	
	$scope.$on('scopeChanged', function() {
		$scope.user = UserService.getUser();
		$scope.currentState = $state.current.name
		console.log($scope.currentState)
	});	
}]);

blogAdminControllers.controller('AdminHomeCtrl', ['$rootScope', '$scope', '$http', '$state', '$timeout', '$sce', 'UserService',
	function($rootScope, $scope, $http, $state, $timeout, $sce, UserService) {

	$scope.user = UserService.getUser();

	$scope.$watch('user', function() {
		if (!(angular.isUndefined($scope.user))) {
			//console.log($scope.user)

			if($scope.user["displayName"] == ""){
				if ($scope.user["role"] == "SiteAdmin") {
					console.log("New Admin")
					$state.go('root.useredit')
				} else if ($scope.user["role"] == "New") {
					console.log("New User")
					$state.go('root.useredit')
				} else {
					console.log("User Loaded")
				}
			}
		} else {
			console.log("Undefined User")
		}
	});
   
	$scope.$on('scopeChanged', function() {
		$scope.user = UserService.getUser();
		$scope.currentState = $state.current.name
		console.log($scope.currentState)
	});
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

	
	$scope.update = function(user) {
	 	$http.post('/api/users', user).success(function(data) {
	        $timeout(function() {
				$scope.user = data;
	        }, 100);
			UserService.saveUser(data);
			$rootScope.$broadcast('scopeChanged', "root.home")
	        //$state.go('root.home', {}, {reload: true});			
			$state.go('root.home')
		})
	};
	
	$scope.cancelEdit = function() {
		$state.go('root.home')
	}
	
	$scope.add = function(newUser) {
	 	$http.post('/api/users', newUser).success(function(data) {
	        $timeout(function() {
				$scope.user = data;
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
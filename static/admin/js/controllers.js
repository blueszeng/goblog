var blogAdminControllers = angular.module('blogAdminControllers', []);

blogAdminControllers.service('UserService', [
	function() {
		
	var currentUser = []

	var getUser = function () {
		console.log("Retrieve User")
		console.log(currentUser[0])
		return currentUser[0];
	};
	
	var saveUser = function(newObj) {
		currentUser.splice(0)
		console.log("Store User")
		currentUser.push(newObj);
	};

	return {
		saveUser: saveUser,
		getUser: getUser
	};
}]);

blogAdminControllers.service('BlogIndexService', ['$rootScope', '$timeout', '$http',
	function($rootScope, $timeout, $http) {
		
	var blogsIndex = []
	var blogIndex = []
		
	var saveBlogs = function(newObj) {
		blogsIndex.splice(0)
		blogsIndex.push(newObj);
	};
	
	var saveBlog = function(newObj) {
		blogIndex.splice(0);
		blogIndex.push(newObj);
	}
	
	var getBlogs = function() {
//		if (!blogsIndex[0]) {
//			saveBlogs();
//		}
		console.log(blogsIndex[0])
		return blogsIndex[0];
	};
	
	var getBlog = function(string) {
		if (!blogIndex[0]) {
			console.log("No Saved Blog")
			$timeout(function() {
				$http.get('/api/blogs/'+string).then(function(data) {
					if (data.error == "No Blogs Found") {
						console.log("No Blog Found")
						return "No Blog Found"
					} else {
						console.log(data.data)
						saveBlog(data.data);
						$rootScope.$broadcast('scopeChanged', "root.posts")
					};										
				});	
			}, 0);
		} else {
			console.log("Saved Blog Found")
			console.log(blogIndex[0]);
			return blogIndex[0];			
		};
	};
	
	var loadBlog = function(string) {
		var $blog = null;
		$http.get('/api/blogs/'+string).then(function(data) {
			if (data.error == "No Blogs Found") {
				console.log("No Blog Found")
				$blog = "No Blog Found";
			} else {
				$blog = data;
			};
		});
		return $blog;
	};
	
	
	return {
		saveBlogs: saveBlogs,
		saveBlog: saveBlog,
		getBlogs: getBlogs,
		getBlog: getBlog,
		loadBlog: loadBlog
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
	};
	
	$scope.add = function(newUser) {
	 	$http.post('/api/users', newUser).success(function(data) {
	        $timeout(function() {
				$scope.user = data;
	            $state.go('^', {}, {reload: true});
	        }, 100);
		})
	};

    	
}]);
	
blogAdminControllers.controller('BlogsListCtrl', ['$rootScope', '$scope', '$http', '$state', '$timeout', 'UserService', 'BlogIndexService',
	function($rootScope, $scope, $http, $state, $timeout, UserService, BlogIndexService) {

	$scope.blogLoad = function () {
		$http.get('/api/blogs/all').success(function(data) {
			if (data.error == "No Blogs Found") {
				console.log("No Blogs")
				$scope.hasBlogs = false
				$scope.loaded = true
			} else {
				$scope.hasBlogs = true
				$scope.blogs = data
				BlogIndexService.saveBlogs(data);
				$scope.loaded = true
			};
		});
	};
	
	$scope.blogLoad();

	$scope.openBlog = function(blog) {
		BlogIndexService.saveBlog(blog);
	     $state.go('root.posts', {blogID: blog.id}, {reload: true});
	}
	
	$rootScope.$broadcast('scopeChanged', "root.blogs")
}]);

blogAdminControllers.controller('BlogEditCtrl', ['$scope', '$http', '$stateParams', '$state', '$timeout', '$location', '$anchorScroll', 'BlogIndexService',
	function($scope, $http, $stateParams, $state, $timeout, $location, $anchorScroll, BlogIndexService) {
    	$scope.blogID = $stateParams.blogID;
		$scope.blogs = BlogIndexService.getBlogs();

		if (!$scope.blogID) {
			console.log("New Blog")
	    	$http.get('/api/blogs/new').success(function(data) {
	    		$scope.blog = data
	    	});			
		} else {
	    	$http.get('/api/blogs/'+$scope.blogID).success(function(data) {
	    		$scope.blog = data
	    	});			
		};

    	$scope.update = function(blog) {
	     	$http.post('/api/blogs', blog).success(function() {
	            $timeout(function() {
	            	$state.go('^', {}, {reload: true});
	            }, 100);
	     	})
      	};
      	
      	$scope.deleteEmail = function (index) {
        	$scope.blog.blogAuthors.splice(index, 1);
    	}
    	
    	$scope.addEmail = function (index) {
    		$http.get('/api/userlookup/' + $scope.newEmail).success(function(data) {
				console.log("data is:", data)
				if (data.Email == "") {
					console.log("No User Found");
					$scope.newEmail = null;
					$scope.emailNotFound = true;
				} else {
        			$scope.blog.blogAuthors.push(data);
					$scope.newEmail = null;
					$scope.emailNotFound = false; 		
    			};
			})
    	}
    	
    	$scope.cancelEdit = function() {
    		$state.go('^')
    	}
    }]);
	
blogAdminControllers.controller('PostsListCtrl', ['$rootScope', '$scope', '$http', '$stateParams', '$state', '$timeout', '$filter', 'UserService', 'BlogIndexService',
	function($rootScope, $scope, $http, $stateParams, $state, $timeout, $filter, UserService, BlogIndexService) {

    $scope.blogID = $stateParams.blogID;

	$scope.$on('scopeChanged', function() {
		$scope.blog = BlogIndexService.getBlog($scope.blogID);
		if ($scope.blog) {
			$scope.blogLoaded = true;
			console.log($scope.blog.sortMethod)
			switch($scope.blog.sortMethod) {
				case '1':
					console.log("Newest Post on top")
					$scope.sort = 'postDate'
					break;
				case '2':
					console.log("Oldest Post on top")
					$scope.sort = '-postDate'
					break;	
				case '3':
					console.log("Custom Order")
					$scope.sort = 'position'
					break;		
			}
			$scope.postsLoad();
			
		};
	});

	$scope.postsLoad = function () {
		$http.get('/api/posts/'+$scope.blogID+'/all').success(function(data) {
			if (data.error == "No Posts Found") {
				console.log("No Posts")
				$scope.hasPosts = false
				$scope.loaded = true
			} else {
				$scope.hasPosts = true
				$scope.posts = data
				$scope.loaded = true
			};
		});
	};
	
}]);
	
blogAdminControllers.controller('PostEditCtrl', ['$scope', '$http', '$stateParams', '$state', '$timeout', '$location', '$anchorScroll', 'BlogIndexService',
	function($scope, $http, $stateParams, $state, $timeout, $location, $anchorScroll, BlogIndexService) {
    	$scope.blogID = $stateParams.blogID;
    	$scope.postID = $stateParams.postID;		
		$scope.blogs = BlogIndexService.getBlogs();

		console.log($scope.blogID)
		console.log($scope.postID)
		
		if (!$scope.postID) {
			console.log("New Blog")
	    	$http.get('/api/posts/'+$scope.blogID+'/new').success(function(data) {
	    		$scope.post = data
	    	});			
		} else {
	    	$http.get('/api/posts/'+$scope.blogID+'/'+$scope.postID).success(function(data) {
	    		$scope.post = data
	    	});			
		};

    	$scope.update = function(post) {
	     	$http.post('/api/posts/'+$scope.blogID, post).success(function() {
	            $timeout(function() {
	            	$state.go('^', {}, {reload: true});
	            }, 100);
	     	})
      	};
      	
    	$scope.cancelEdit = function() {
    		$state.go('^')
    	};

}]);
	
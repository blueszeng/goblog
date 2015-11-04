var blogAdminControllers = angular.module('blogAdminControllers', []);

var debugFlag = true

blogAdminControllers.service('UserService', ['$rootScope', '$timeout', '$http',
	function($rootScope, $timeout, $http) {
		
	var currentUser = []

	var getUser = function () {
		if(debugFlag){console.log("UserService.getUser Entered")};
		if (!currentUser[0]) {
			if(debugFlag){console.log("UserService.getUser No Stored User")};
			$timeout(function() {
				$http.get('/api/users').success(function(data) {
					if (data.error == "No Users Found") {
						if(debugFlag){console.log("UserService.getUser /api/users No User Found")}
						return "No User Found"
					} else {
						if(debugFlag){console.log("UserService.getUser /api/users User Found")}
						saveUser(data);
						console.log(data.role)
						$rootScope.$broadcast('scopeChanged', "root.home")
						//return(data);
					};										
				});
			});
		} else {
			if(debugFlag) {console.log("UserService.getUser Retrieve Stored User")};
			return currentUser[0];			
		}
	};

	var retrieveUser = function() {
		if(debugFlag){console.log("UserService.retrieveUser Entered")};
		return currentUser[0];			
	};
	
	var saveUser = function(newObj) {
		if(debugFlag){console.log("UserService.saveUser Entered")};
		currentUser.splice(0)
		currentUser.push(newObj);
	};

	return {
		saveUser: saveUser,
		getUser: getUser,
		retrieveUser: retrieveUser
	};
}]);

blogAdminControllers.service('BlogIndexService', ['$rootScope', '$timeout', '$http',
	function($rootScope, $timeout, $http) {
		
	var blogIndex = []
	var postIndex = []

//	var blogsIndex = []	
//	var postsIndex = []

//	var saveBlogs = function(newObj) {
//		blogsIndex.splice(0)
//		blogsIndex.push(newObj);
//	};

//	var getBlogs = function() {
////		if (!blogsIndex[0]) {
////			saveBlogs();
////		}
//		console.log(blogsIndex[0])
//		return blogsIndex[0];
//	};
	
	var saveBlog = function(newObj) {
		if(debugFlag){console.log("BlogIndexService.saveBlog Entered")};
		blogIndex.splice(0);
		blogIndex.push(newObj);
	};

	var savePost = function(newObj) {
		if(debugFlag){console.log("BlogIndexService.savePost Entered")};
		postIndex.splice(0);
		postIndex.push(newObj);
	};
		
	var getBlog = function(string) {
		if(debugFlag){console.log("BlogIndexService.getBlog Entered")};
		if (!blogIndex[0]) {
			if(debugFlag){console.log("BlogIndexService.getBlog No Saved Blog")};
			$timeout(function() {
				$http.get('/api/blogs/'+string).then(function(data) {
					if (data.error == "No Blogs Found") {
						if(debugFlag){console.log("BlogIndexService.getBlog /api/blogs/ No Blog Found")};
						return "No Blog Found"
					} else {
						if(debugFlag){console.log("BlogIndexService.getBlog /api/blogs/ Blog Found")};
						saveBlog(data.data);
						console.log(data.data);
						//return data.data;
						$rootScope.$broadcast('blogLoaded', "root.posts")
					};										
				});	
			}, 0);
		} else {
			if(debugFlag){console.log("BlogIndexService.getBlog Retrieve Stored Blog")};			
			return blogIndex[0];			
		};
	};

	var getPost = function(blogID, postID) {
		if(debugFlag){console.log("BlogIndexService.getPost Entered")};
		if (!postIndex[0]) {			
			if(debugFlag){console.log("BlogIndexService.getPost No Saved Post")};
			$timeout(function() {
				$http.get('/api/posts/'+blogID+'/'+postID).then(function(data) {				
					if (data.error == "No Posts Found") {
						if(debugFlag){console.log("BlogIndexService.getPost /api/posts/ No Post Found")};
						return "No Post Found"
					} else {
						if(debugFlag){console.log("BlogIndexService.getPost /api/posts/ Post Found")};
						console.log(data.data)
						savePost(data.data);
						$rootScope.$broadcast('postLoaded', "root.entries")
					};										
				});	
			}, 0);
		} else {
			if(debugFlag){console.log("BlogIndexService.getPost Retrieve Stored Post")};			
			return postIndex[0];			
		};
	};
	
//	var loadBlog = function(string) {
//		var $blog = null;
//		$http.get('/api/blogs/'+string).then(function(data) {
//			if (data.error == "No Blogs Found") {
//				console.log("No Blog Found")
//				$blog = "No Blog Found";
//			} else {
//				$blog = data;
//			};
//		});
//		return $blog;
//	};
	
	return {
//		saveBlogs: saveBlogs,
//		getBlogs: getBlogs,
//		loadBlog: loadBlog,
		saveBlog: saveBlog,
		getBlog: getBlog,
		savePost: savePost,
		getPost: getPost
	};
}]);

blogAdminControllers.controller('AdminHeaderCtrl', ['$rootScope', '$scope', '$http', '$timeout', '$state', 'UserService',
	function($rootScope, $scope, $http, $timeout, $state, UserService) {
	if(debugFlag){console.log("AdminHeaderCtrl Entered")};	
	
	$scope.user = UserService.getUser();
	
	$scope.$on('scopeChanged', function() {
		if(debugFlag){console.log("AdminHeaderCtrl.scopeChanged Retrieving User")};	
		//$scope.user = UserService.getUser();
		$scope.user = UserService.retrieveUser();

		$scope.currentState = $state.current.name
	});	
}]);

blogAdminControllers.controller('AdminHomeCtrl', ['$rootScope', '$scope', '$http', '$state', '$timeout', '$sce', 'UserService',
	function($rootScope, $scope, $http, $state, $timeout, $sce, UserService) {
	if(debugFlag){console.log("AdminHomeCtrl Entered")};	

	$scope.user = UserService.retrieveUser();

	$scope.$watch('user', function() {
		if(debugFlag){console.log("AdminHomeCtrl.watch(user)")};	
		if (!(angular.isUndefined($scope.user))) {
			if($scope.user["displayName"] == ""){
				if ($scope.user["role"] == "SiteAdmin") {
					$state.go('root.useredit')
				} else if ($scope.user["role"] == "New") {
					$state.go('root.useredit')
				} else {
					console.log("AdminHomeCtrl.watch(user) User Defined")
				}
			}
		} else {
			if(debugFlag){console.log("AdminHomeCtrl.watch(user) User Undefined")};	
		}
	});
   
	$scope.$on('scopeChanged', function() {
		if(debugFlag){console.log("AdminHomeCtrl.scopeChanged Retrieving User")};	
		$scope.user = UserService.retrieveUser();
		$scope.currentState = $state.current.name
	});
}]);

blogAdminControllers.controller('UsersListCtrl', ['$rootScope', '$scope', '$http', '$state', '$timeout', 'UserService',
	function($rootScope, $scope, $http, $state, $timeout, UserService) {
	if(debugFlag){console.log("UserListCtrl Entered")};	

	$scope.user = UserService.retrieveUser()

	$scope.$watch('user', function() {
		if(debugFlag){console.log("UsersListCtrl.watch(user)")};	
		if (angular.isUndefined($scope.user)) {
			if(debugFlag){console.log("UsersListCtrl.watch(user) User Undefined")};	
		} else if ($scope.user["role"] != "SiteAdmin") {
			if(debugFlag){console.log("UsersListCtrl.watch(user) User Not Authorized")};	
			$state.go('root.home')
		}
	});
   
	$http.get('/api/users/all').success(function(data) {
		if(debugFlag){console.log("UsersListCtrl.get Load All Users")};	
		$timeout(function() {
			$scope.users = data
		}, 0);
	});
	
	if(debugFlag){console.log("UsersListCtrl.broadcast root.users scopeChange")};	
	$rootScope.$broadcast('scopeChanged', "root.users")
}]);

blogAdminControllers.controller('UserEditCtrl', ['$rootScope', '$scope', '$http', '$state', '$timeout', 'UserService',
	function($rootScope, $scope, $http, $state, $timeout, UserService) {
	if(debugFlag){console.log("UserEditCtrl Entered")};	

	$scope.user = UserService.retrieveUser()

	$scope.$watch('user', function() {
		if(debugFlag){console.log("UserEditCtrl.watch(user)")};	
		if (angular.isUndefined($scope.user)) {
			if(debugFlag){console.log("UserEditCtrl.watch(user) User Undefined")};	
			//console.log("No User Logged In")
			//$scope.user = UserService.retrieveUser()
			$state.go('root.home')	   	
		} else if ($scope.user["email"] == "") {
			if(debugFlag){console.log("UserEditCtrl.watch(user) User Email Missing")};	
			//console.log("No User Email Address")
			$state.go('root.home')
		} else {
			if(debugFlag){console.log("UserEditCtrl.watch(user) User Found")};	
			//console.log($scope.user["role"])
		}
	});

	$scope.update = function(user) {
		if(debugFlag){console.log("UserEditCtrl.update Entered")};	
	 	if (user.displayName == "") {
			if(debugFlag){console.log("UserEditCtrl.update Name Not Found")};	
			$scope.nameNotFound = true;
		} else {
			$scope.nameNotFound = false;
			if(debugFlag){console.log("UserEditCtrl.update Post User")};	
			$http.post('/api/users', user).success(function(data) {
		        $timeout(function() {
					$scope.user = data;
		        }, 100);
				UserService.saveUser(data);
				$rootScope.$broadcast('scopeChanged', "root.home")
		        $state.go('root.home', {}, {reload: true});			
			});			
		}
	};
	
	$scope.cancelEdit = function() {
		if(debugFlag){console.log("UserEditCtrl.cancelEdit")};	
		$state.go('root.home')
	};
	
	$scope.add = function(newUser) {
		if(debugFlag){console.log("UserEditCtrl.add Entered")};	
		if (newUser == null) {
			if(debugFlag){console.log("UserEditCtrl.add Email Not Found")};	
			$scope.emailNotFound = true;
		} else if (!newUser.email.match(/.+@.+\..+/i)) {
			if(debugFlag){console.log("UserEditCtrl.add Email Not Valid")};	
			$scope.emailNotFound = true;
			$scope.newUser.email = null;
		} else {					
			$scope.emailNotFound = false;
			if(debugFlag){console.log("UserEditCtrl.add Post User")};	
		 	$http.post('/api/users', newUser).success(function(data) {
		        $timeout(function() {
					$scope.user = data;
		            $state.go('^', {}, {reload: true});
		        }, 100);
			});			
		};
	};

    	
}]);
	
blogAdminControllers.controller('BlogsListCtrl', ['$rootScope', '$scope', '$http', '$state', '$timeout', 'UserService', 'BlogIndexService',
	function($rootScope, $scope, $http, $state, $timeout, UserService, BlogIndexService) {
	if(debugFlag){console.log("BlogsListCtrl Entered")};	

	$scope.user = UserService.retrieveUser()

	$scope.$watch('user', function() {
		if(debugFlag){console.log("BlogsListCtrl.watch(user)")};	
		if (angular.isUndefined($scope.user)) {
			if(debugFlag){console.log("BlogsListCtrl.watch(user) User Undefined")};	
			//$state.go('root.home')	   	
		} else if ($scope.user["role"] != "SiteAdmin") {
			if(debugFlag){console.log("BlogsListCtrl.watch(user) User Not Authorized")};	
			$state.go('root.home')
		}
	});

	$scope.blogLoad = function () {
		if(debugFlag){console.log("BlogsListCtrl.blogLoad Entered")};	
		$http.get('/api/blogs/all').success(function(data) {
			if (data.error == "No Blogs Found") {
				if(debugFlag){console.log("BlogsListCtrl.blogLoad No Blogs Found")};	
				console.log("No Blogs")
				$scope.hasBlogs = false
				$scope.loaded = true
			} else {
				if(debugFlag){console.log("BlogsListCtrl.blogLoad Blogs Loaded")};	
				$scope.hasBlogs = true
				$scope.blogs = data
				$scope.loaded = true
			};
		});
	};
	
	$scope.blogLoad();

	$scope.openBlog = function(blog) {
		if(debugFlag){console.log("BlogsListCtrl.openBlog Entered")};	
		BlogIndexService.saveBlog(blog);
	    $state.go('root.posts', {blogID: blog.id}, {reload: true});
	}

	if(debugFlag){console.log("BlogssListCtrl.broadcast root.blogs scopeChange")};		
	$rootScope.$broadcast('scopeChanged', "root.blogs")
}]);

blogAdminControllers.controller('BlogEditCtrl', ['$scope', '$http', '$stateParams', '$state', '$timeout', '$location', '$anchorScroll', 'BlogIndexService',
	function($scope, $http, $stateParams, $state, $timeout, $location, $anchorScroll, BlogIndexService) {
		if(debugFlag){console.log("BlogsEditCtrl Entered")};	

    	$scope.blogID = $stateParams.blogID;

		if (!$scope.blogID) {
			if(debugFlag){console.log("BlogsEditCtrl New Blog")};	
	    	$http.get('/api/blogs/new').success(function(data) {
	    		$scope.blog = data
	    	});			
		} else {
			if(debugFlag){console.log("BlogsEditCtrl Load Blog")};				
	    	$http.get('/api/blogs/'+$scope.blogID).success(function(data) {
	    		$scope.blog = data
	    	});			
		};

    	$scope.update = function(blog) {
			if(debugFlag){console.log("BlogsEditCtrl.update Entered")};				
    		if (blog.blogName == "") {
				if(debugFlag){console.log("BlogsEditCtrl.update Title Not Entered")};				
    			$scope.titleNotFound = true
    		} else {
    			$scope.titleNotFound = false
			if(debugFlag){console.log("BlogsEditCtrl.update Blog Post")};								
		     	$http.post('/api/blogs', blog).success(function() {
		            $timeout(function() {
		            	$state.go('^', {}, {reload: true});
		            }, 100);
		     	})
    		};
      	};
      	
      	$scope.deleteEmail = function (index) {
			if(debugFlag){console.log("BlogsEditCtrl.deleteEmail Entered")};				
        	$scope.blog.blogAuthors.splice(index, 1);
    	}
    	
    	$scope.addEmail = function (index) {
			if(debugFlag){console.log("BlogsEditCtrl.addEmail Entered")};							
			if ($scope.newEmail == null) {
				if(debugFlag){console.log("BlogsEditCtrl.addEmail Email Not Entered")};								
				$scope.emailNotFound = true;
			} else if (!$scope.newEmail.match(/.+@.+\..+/i)) {
				if(debugFlag){console.log("BlogsEditCtrl.addEmail Email Not Valid")};												
				$scope.emailNotFound = true;
				$scope.newEmail = null;
			} else {
				if(debugFlag){console.log("BlogsEditCtrl.addEmail Email Lookup")};								
	    		$http.get('/api/userlookup/' + $scope.newEmail).success(function(data) {
					console.log("data is:", data)
					if (data.Email == "") {
						if(debugFlag){console.log("BlogsEditCtrl.addEmail Email Not Found")};								
						console.log("No User Found");
						$scope.newEmail = null;
						$scope.emailNotFound = true;
					} else {
						if(debugFlag){console.log("BlogsEditCtrl.addEmail Email Lookup")};								
	        			$scope.blog.blogAuthors.push(data);
						$scope.newEmail = null;
						$scope.emailNotFound = false; 		
	    			};
				})
			}
    	}
    	
    	$scope.cancelEdit = function() {
			if(debugFlag){console.log("BlogsEditCtrl.cancelEdit Entered")};
    		$state.go('^')
    	}
    }]);
	
blogAdminControllers.controller('PostsListCtrl', ['$rootScope', '$scope', '$http', '$stateParams', '$state', '$timeout', '$filter', 'UserService', 'BlogIndexService',
	function($rootScope, $scope, $http, $stateParams, $state, $timeout, $filter, UserService, BlogIndexService) {
	if(debugFlag){console.log("PostsListCtrl Entered")};	

    $scope.blogID = $stateParams.blogID;

	$scope.user = UserService.retrieveUser()
	$scope.blog = BlogIndexService.getBlog($scope.blogID);

	$scope.$watch('user', function() {
		if (angular.isUndefined($scope.user)) {
			if(debugFlag){console.log("PostsListCtrl.watch(user) User Undefined")};	
			//$state.go('root.home')	   	
		} else if ($scope.user["role"] != "SiteAdmin") {
			if(debugFlag){console.log("PostsListCtrl.watch(user) User Not Authorized")};	
			$state.go('root.home')
		}
	});
	
	$scope.$watch('blog', function() {
		if ($scope.blog) {
			$scope.blogLoaded = true;
			if(debugFlag){console.log("PostsListCtrl.watch(blog) Blog Founded")};				
			//console.log($scope.blog.sortMethod)
			switch($scope.blog.sortMethod) {
				case '1':
					console.log("Newest Post on top")
					$scope.sort = '-postDate'
					break;
				case '2':
					console.log("Oldest Post on top")
					$scope.sort = 'postDate'
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
		if(debugFlag){console.log("PostsListCtrl.postsLoad Entered")};			
		$http.get('/api/posts/'+$scope.blogID+'/all').success(function(data) {
			if (data.error == "No Posts Found") {
				if(debugFlag){console.log("PostsListCtrl.postsLoad No Posts Found")};			
				console.log("No Posts")
				$scope.hasPosts = false
				$scope.loaded = true
			} else {
		if(debugFlag){console.log("PostsListCtrl.postsLoad Posts Loaded")};							
				$scope.hasPosts = true
				$scope.posts = data
				$scope.loaded = true
			};
		});
	};

//	$scope.blogLoad = function() {
//		if(debugFlag){console.log("PostsListCtrl.blogLoad Entered")};	
//		$scope.blog = BlogIndexService.getBlog($scope.blogID);
//		console.log($scope.blog)
//	};

	//$scope.blogLoad();		
	
	$scope.$on('blogLoaded', function() {
		if(debugFlag){console.log("PostsListCtrl.on scopeChange")};	
		$scope.blog = BlogIndexService.getBlog($scope.blogID);
	});
	
	$scope.openPost = function(post) {
		if(debugFlag){console.log("PostsListCtrl.openPost Entered")};					
    	blog_ID = $stateParams.blogID;
		BlogIndexService.savePost(post);
	     $state.go('root.entries', {blogID: blog_ID, postID: post.id}, {reload: true});
	}

	if(debugFlag){console.log("PostsListCtrl.broadcast root.posts scopeChange")};		
	$rootScope.$broadcast('scopeChanged', "root.posts")	
}]);
	
blogAdminControllers.controller('PostEditCtrl', ['$scope', '$http', '$stateParams', '$state', '$timeout', '$location', '$anchorScroll', 'BlogIndexService',
	function($scope, $http, $stateParams, $state, $timeout, $location, $anchorScroll, BlogIndexService) {
		if(debugFlag){console.log("PostEditCtrl Entered")};	
	 	$scope.blogID = $stateParams.blogID;
    	$scope.postID = $stateParams.postID;		

		console.log($scope.blogID)
		console.log($scope.postID)
		
		if (!$scope.postID) {
			if(debugFlag){console.log("PostEditCtrl New Post")};	
	    	$http.get('/api/posts/'+$scope.blogID+'/new').success(function(data) {
	    		$scope.post = data
	    	});			
		} else {
			if(debugFlag){console.log("PostEditCtrl Load Post")};	
	    	$http.get('/api/posts/'+$scope.blogID+'/'+$scope.postID).success(function(data) {
	    		$scope.post = data
	    	});			
		};

    	$scope.update = function(post) {
			if(debugFlag){console.log("PostEditCtrl.update Entered")};				
			if (post.postName == "") {
				if(debugFlag){console.log("PostEditCtrl.update Title Not Entered")};				
    			$scope.titleNotFound = true
    		} else {
				$scope.titleNotFound = false
			}

			if (post.postDateStr == "" || !post.postDateStr.match(/^[0-3]?[0-9]-[0-3]?[0-9]-[0-9]{4}$/)) {
				if(debugFlag){console.log("PostEditCtrl.update Post Date Not Entered")};				
				$scope.postDateNotFound = true
			} else {
				$scope.postDateNotFound = false
			}
			
			if ((post.stopDateStr == "" || !post.stopDateStr.match(/^[0-3]?[0-9]-[0-3]?[0-9]-[0-9]{4}$/)) && post.stopFlag) {
				if(debugFlag){console.log("PostEditCtrl.update Stop Date Not Entered")};				
				$scope.stopDateNotFound = true
			} else {
				$scope.stopDateNotFound = false
			}
			
			if (!$scope.titleNotFound && !$scope.postDateNotFound && !$scope.stopDateNotFound) {
				if(debugFlag){console.log("PostEditCtrl.update Post Post")};								
		     	$http.post('/api/posts/'+$scope.blogID, post).success(function() {
		            $timeout(function() {
		            	$state.go('^', {}, {reload: true});
		            }, 100);
		     	});
			};
      	};
      	
    	$scope.cancelEdit = function() {
    		$state.go('^')
    	};

}]);

blogAdminControllers.controller('EntryCtrl', ['$rootScope', '$scope', '$http', '$stateParams', '$state', '$timeout', '$filter', 'UserService', 'BlogIndexService',
	function($rootScope, $scope, $http, $stateParams, $state, $timeout, $filter, UserService, BlogIndexService) {
	if(debugFlag){console.log("EntryCtrl Entered")};	

    $scope.blogID = $stateParams.blogID;
    $scope.postID = $stateParams.postID;

	$scope.user = UserService.retrieveUser()

	$scope.$watch('user', function() {
		if (angular.isUndefined($scope.user)) {
			if(debugFlag){console.log("EntryCtrl.watch(user) User Undefined")};	
			//$state.go('root.home')	   	
		} else if ($scope.user["role"] != "SiteAdmin") {
			if(debugFlag){console.log("EntryCtrl.watch(user) User Not Authorized")};	
			$state.go('root.home')
		}
	});

	$scope.post = BlogIndexService.getPost($scope.blogID, $scope.postID);

	$scope.$on('postLoaded', function() {
		if(debugFlag){console.log("EntryCtrl.on postLoaded scopeChange")};	
		$scope.post = BlogIndexService.getPost($scope.blogID, $scope.postID);
	});

//	$scope.$on('scopeChanged', function() {
//		$scope.post = BlogIndexService.getPost($scope.blogID, $scope.postID);
//		if ($scope.post) {
//			$scope.postLoaded = true;
//			console.log($scope.post)
//			//$scope.entriesLoad();	
//		};
//	});

	$scope.$watch('post', function() {
		if ($scope.post) {
			$scope.postLoaded = true;
			if(debugFlag){console.log("EntryCtrl.watch(post) Post Founded")};				
			$scope.entriesLoad();
		};		
	});	


	$scope.entriesLoad = function () {
		$http.get('/api/entries/'+$scope.postID+'/latest').success(function(data) {
			if (data.error == "No Entries Found") {
				console.log("No Entry")
				$scope.hasEntry = false
				$scope.loaded = true
				$state.go('root.entries.edit')
			} else {
				$scope.hasEntry = true
				$scope.entry = data
				$scope.loaded = true
			};
		});
	};
	
	if(debugFlag){console.log("EntryCtrl.broadcast root.entries scopeChange")};		
	$rootScope.$broadcast('scopeChanged', "root.entries")	
}]);
	
blogAdminControllers.controller('EntryEditCtrl',  ['$rootScope', '$scope', '$http', '$stateParams', '$state', '$timeout', 'BlogIndexService',
	function($rootScope, $scope, $http, $stateParams, $state, $timeout, BlogIndexService) {
		if(debugFlag){console.log("EntryEditCtrl Entered")};	
	 	$scope.blogID = $stateParams.blogID;
    	$scope.postID = $stateParams.postID;		
		$scope.entryID = $stateParams.entryID;
		
		console.log($scope.blogID)
		console.log($scope.postID)
		console.log($scope.entryID)

		if (!$scope.entryID) {
			if(debugFlag){console.log("EntryEditCtrl New Entry")};	
	    	$http.get('/api/entries/'+$scope.postID+'/new').success(function(data) {
	    		$scope.entry = data;
	    	});			
		} else {
			if(debugFlag){console.log("EntryEditCtrl Load Post")};	
	    	$http.get('/api/entries/'+$scope.postID+'/'+$scope.entryID).success(function(data) {
	    		$scope.entry = data;
	    	});			
		};
		
    	$scope.update = function(entry) {
			if(debugFlag){console.log("EntryEditCtrl.update Entered")};				
			if(debugFlag){console.log("EntryEditCtrl.update Entry Post")};								
	     	$http.post('/api/entries/'+$scope.postID, entry).success(function() {
	            $timeout(function() {
	            	$state.go('^', {}, {reload: true});
	            }, 100);
	     	});
      	};
      	
    	$scope.cancelEdit = function() {
    		$state.go('^')
    	};
}]);
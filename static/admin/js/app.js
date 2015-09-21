var blogAdminApp = angular.module('blogAdminApp', [
  'ui.router',
  'mm.foundation',
  'blogAdminControllers',
  'myDirectiveModule'
]);

/*
blogAdminApp.config(function($stateProvider) {
  $stateProvider.state("Modal", {
    views:{
      "modal": {
        templateUrl: "index.html",
        controller: 'AdminHeaderCtrl'
      }
    },
    abstract: true
  });
});
*/

blogAdminApp.config(['$stateProvider', '$urlRouterProvider', function($stateProvider, $urlRouterProvider) {
    $urlRouterProvider.otherwise('/');
 
    $stateProvider
		.state('root', {
		  url:'',
		  abstract: true,
		  views: {
		    'header': {
		      templateUrl: "templates/header.html",
		      controller: 'AdminHeaderCtrl'  
		    }
		  }
		})
		.state('root.home', {
		  url:'/',
		  views: {
		    'container@': {
		      templateUrl:"templates/home.html",
		      controller: 'AdminHomeCtrl'
		    }
		  }
		})
		.state('root.users', {
		  url:'/users',
		  views: {
		    'container@': {
		      templateUrl: 'templates/users.html',
		      controller: 'UsersListCtrl'
		    }
		  }
		})
		.state('root.useredit', {
		  url:'/user/edit',
		  views: {
		    'container@': {
		      templateUrl: 'templates/userEdit.html',
		      controller: 'UserEditCtrl'
		    }
		  }
		})
		.state('root.users.add', {
		  url:'/add',
		  views: {
		    'subcontainer': {
		      templateUrl: 'templates/userAdd.html',
		      controller: 'UserEditCtrl'
		    }
		  }
		})
		.state('root.blogs', {
		  url:'/blogs',
		  views: {
		    'container@': {
		      templateUrl: 'templates/blogs.html',
		      controller: 'BlogsListCtrl'
		    }
		  }
		})  
		.state('root.blogs.edit', {
		  url:'/:blogID',
		  views: {
		    'subcontainer': {
		      templateUrl: 'templates/blogEdit.html',
		      controller: 'BlogEditCtrl'
		    }
		  }
		})
		.state('root.posts', {
			url:'/blog/:blogID/posts',
			views: {
				'container@': {
				templateUrl: 'templates/posts.html',
				controller: 'PostsListCtrl'
				}
			}
		})
		.state('root.posts.edit', {
			url:'/:postID',
			views: {
				'subcontainer': {
				templateUrl: 'templates/postEdit.html',
				controller: 'PostEditCtrl'
				}
			}
		})	 
}]);

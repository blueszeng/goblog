var blogAdminServices = angular.module('blogAdminServices', []);

blogAdminServices.service('userService', ['$scope', '$http', '$timeout',
	function($scope, $http, $timeout) {
		
  var currentUser = {
  	email: "none",
  	displayName : "none",
  	role: "guest"
  }

  this.getUser = function (){
      return currentUser;
  };

	this.loadUser = function (){
		$http.get('/api/users').success(function(data) {
		$timeout(function() {
			currentUser = data;
		}, 50);
		});
	};

}]);
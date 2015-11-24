var openinstrumentApp = angular.module('openinstrumentApp', ['ngMaterial']);

openinstrumentApp.controller('DatastoreStatusCtrl', function($scope, $http, $interval, $mdSidenav) {
	$scope.loadBlocks = function() {
		$http.get('/blocks').success(function(data) {
			$scope.blocks = data;
		});
	};

	$scope.loadBlocks();
	$interval($scope.loadBlocks, 5000);
});

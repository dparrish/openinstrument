var openinstrumentApp = angular.module('openinstrumentApp', ['ngMaterial']);

openinstrumentApp.controller('DatastoreStatusCtrl', function($scope, $http, $interval, $mdSidenav) {
	$scope.total_values = 0;
	$scope.total_streams = 0;

	$scope.loadBlocks = function() {
		$http.get('/blocks').success(function(data) {
			$scope.blocks = data;

			$scope.indexed_values = 0;
			$scope.indexed_streams = 0;
			$scope.logged_values = 0;
			$scope.logged_streams = 0;
			$scope.unlogged_values = 0;
			$scope.unlogged_streams = 0;
			$scope.total_size = 0;

			angular.forEach(data,
				function(obj, key) {
					$scope.indexed_streams += obj.indexed_streams || 0;
					$scope.indexed_values += obj.indexed_values || 0;
					$scope.logged_streams += obj.logged_streams || 0;
					$scope.logged_values += obj.logged_values || 0;
					$scope.unlogged_streams += obj.unlogged_streams || 0;
					$scope.unlogged_values += obj.unlogged_values || 0;
					$scope.total_size += obj.size || 0;
				}
			)
		});
	};

	$scope.loadBlocks();
	$interval($scope.loadBlocks, 5000);
})
.filter('bytes', function() {
	return function(bytes, precision) {
		if (isNaN(parseFloat(bytes)) || !isFinite(bytes)) return '-';
		if (typeof precision === 'undefined') precision = 1;
		var units = ['bytes', 'kB', 'MB', 'GB', 'TB', 'PB'],
			number = Math.floor(Math.log(bytes) / Math.log(1024));
		return (bytes / Math.pow(1024, Math.floor(number))).toFixed(precision) +  ' ' + units[number];
	}
});

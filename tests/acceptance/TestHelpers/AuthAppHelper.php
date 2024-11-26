<?php declare(strict_types=1);
/**
 * ownCloud
 *
<<<<<<< HEAD
 * @author Sajan Gurung <sajan@jankaritech.com>
 * @copyright Copyright (c) 2024 Sajan Gurung sajan@jankaritech.com
=======
>>>>>>> 8a06712ccc (add test for creating auth tocken for an app using api)
 *
 * This code is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License,
 * as published by the Free Software Foundation;
 * either version 3 of the License, or any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program. If not, see <http://www.gnu.org/licenses/>
 *
 */

namespace TestHelpers;

use Psr\Http\Message\ResponseInterface;

/**
 * A helper class for managing Auth App API requests
 */
class AuthAppHelper {

	/**
	 * @return void
	 */
	public static function getAuthAppEndpoint() {
		return "/auth-app/tokens";
	}

	/**
	 * @param string $baseUrl
	 * @param string $user
	 * @param string $password
	 *
	 * @return ResponseInterface
	 */
<<<<<<< HEAD
	public static function listAllAppAuthToken(string $baseUrl, string $user, string $password) : ResponseInterface {
=======
	public static function listAllAppAuthToken(string $baseUrl,string $user,string $password) : ResponseInterface {
>>>>>>> 8a06712ccc (add test for creating auth tocken for an app using api)
		$url = $baseUrl . self::getAuthAppEndpoint();
		return HttpRequestHelper::sendRequest(
			$url,
			null,
			"GET",
			$user,
			$password,
		);
	}

	/**
	 * @param string $baseUrl
	 * @param string $user
	 * @param string $password
	 * @param string $expiration
	 *
	 * @return ResponseInterface
	 */
<<<<<<< HEAD
	public static function createAppAuthToken(string $baseUrl, string $user, string $password, string $expiration) : ResponseInterface {
=======
	public static function createAppAuthToken(string $baseUrl,string $user,string $password,string $expiration) : ResponseInterface {
>>>>>>> 8a06712ccc (add test for creating auth tocken for an app using api)
		$url = $baseUrl . self::getAuthAppEndpoint() . "?expiry=$expiration";
		return HttpRequestHelper::sendRequest(
			$url,
			null,
			"POST",
			$user,
			$password,
		);
	}

	/**
	 * @param string $baseUrl
	 * @param string $user
	 * @param string $password
	 * @param string $token
	 *
	 * @return ResponseInterface
	 */
	public static function deleteAppAuthToken(string $baseUrl, string $user, string $password, string $token) : ResponseInterface {
		$url = $baseUrl . self::getAuthAppEndpoint() . "?token=$token";
		return HttpRequestHelper::sendRequest(
			$url,
			null,
			"DELETE",
			$user,
			$password,
		);
	}
}

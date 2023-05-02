<?php

declare(strict_types=1);

/**
 * ownCloud
 *
 * @author Viktor Scharf <v.scharf@owncloud.com>
 * @copyright Copyright (c) 2022 Viktor Scharf v.scharf@owncloud.com
 */

use Behat\Behat\Context\Context;
use Behat\Behat\Hook\Scope\BeforeScenarioScope;
use GuzzleHttp\Exception\GuzzleException;
use Behat\Gherkin\Node\TableNode;

require_once 'bootstrap.php';

/**
 * Context for the TUS-specific steps using the Graph API
 */
class SpacesTUSContext implements Context {
	private FeatureContext $featureContext;
	private TUSContext $tusContext;
	private SpacesContext $spacesContext;

	/**
	 * This will run before EVERY scenario.
	 * It will set the properties for this object.
	 *
	 * @BeforeScenario
	 *
	 * @param BeforeScenarioScope $scope
	 *
	 * @return void
	 */
	public function before(BeforeScenarioScope $scope): void {
		// Get the environment
		$environment = $scope->getEnvironment();
		// Get all the contexts you need in this context from here
		$this->featureContext = $environment->getContext('FeatureContext');
		$this->spacesContext = $environment->getContext('SpacesContext');
		$this->tusContext = $environment->getContext('TUSContext');
	}

	/**
	 * @Given /^user "([^"]*)" has uploaded a file from "([^"]*)" to "([^"]*)" via TUS inside of the space "([^"]*)" using the WebDAV API$/
	 *
	 * @param string $user
	 * @param string $source
	 * @param string $destination
	 * @param string $spaceName
	 *
	 * @return void
	 *
	 * @throws Exception
	 * @throws GuzzleException
	 */
	public function userHasUploadedFileViaTusInSpace(string $user, string $source, string $destination, string $spaceName): void {
		$this->userUploadsAFileViaTusInsideOfTheSpaceUsingTheWebdavApi($user, $source, $destination, $spaceName);
		$this->featureContext->theHTTPStatusCodeShouldBe(200, "Expected response status code should be 200");
	}

	/**
	 * @When /^user "([^"]*)" uploads a file from "([^"]*)" to "([^"]*)" via TUS inside of the space "([^"]*)" using the WebDAV API$/
	 *
	 * @param string $user
	 * @param string $source
	 * @param string $destination
	 * @param string $spaceName
	 *
	 * @return void
	 * @throws Exception
	 * @throws GuzzleException
	 */
	public function userUploadsAFileViaTusInsideOfTheSpaceUsingTheWebdavApi(
		string $user,
		string $source,
		string $destination,
		string $spaceName
	): void {
		$this->spacesContext->setSpaceIDByName($user, $spaceName);
		$this->tusContext->userUploadsUsingTusAFileTo($user, $source, $destination);
	}

	/**
	 * @Given user :user has created a new TUS resource for the space :spaceName with content :content using the WebDAV API with these headers:
	 *
	 * @param string $user
	 * @param string $spaceName
	 * @param string $content
	 * @param TableNode $headers
	 *
	 * @return void
	 *
	 * @throws Exception
	 * @throws GuzzleException
	 */
	public function userHasCreatedANewTusResourceForTheSpaceUsingTheWebdavApiWithTheseHeaders(
		string $user,
		string $spaceName,
		string $content,
		TableNode $headers
	): void {
		$this->userCreatesANewTusResourceForTheSpaceUsingTheWebdavApiWithTheseHeaders($user, $spaceName, $content, $headers);
		$this->featureContext->theHTTPStatusCodeShouldBe(201, "Expected response status code should be 201");
	}

	/**
	 * @When user :user creates a new TUS resource for the space :spaceName with content :content using the WebDAV API with these headers:
	 *
	 * @param string $user
	 * @param string $spaceName
	 * @param string $content
	 * @param TableNode $headers
	 *
	 * @return void
	 *
	 * @throws Exception
	 * @throws GuzzleException
	 */
	public function userCreatesANewTusResourceForTheSpaceUsingTheWebdavApiWithTheseHeaders(
		string $user,
		string $spaceName,
		string $content,
		TableNode $headers
	): void {
		$this->spacesContext->setSpaceIDByName($user, $spaceName);
		$this->tusContext->createNewTUSResourceWithHeaders($user, $headers, $content);
	}

	/**
	 * @When /^user "([^"]*)" uploads a file with content "([^"]*)" to "([^"]*)" via TUS inside of the space "([^"]*)" using the WebDAV API$/
	 *
	 * @param string $user
	 * @param string $content
	 * @param string $resource
	 * @param string $spaceName
	 *
	 * @return void
	 * @throws Exception|GuzzleException
	 */
	public function userUploadsAFileWithContentToViaTusInsideOfTheSpaceUsingTheWebdavApi(
		string $user,
		string $content,
		string $resource,
		string $spaceName
	): void {
		$this->spacesContext->setSpaceIDByName($user, $spaceName);
		$this->tusContext->userUploadsAFileWithContentToUsingTus($user, $content, $resource);
	}

	/**
	 * @When /^user "([^"]*)" uploads a file "([^"]*)" to "([^"]*)" with mtime "([^"]*)" via TUS inside of the space "([^"]*)" using the WebDAV API$/
	 *
	 * @param string $user
	 * @param string $source
	 * @param string $destination
	 * @param string $mtime Time in human-readable format is taken as input which is converted into milliseconds that is used by API
	 * @param string $spaceName
	 *
	 * @return void
	 *
	 * @throws Exception
	 * @throws GuzzleException
	 */
	public function userUploadsAFileToWithMtimeViaTusInsideOfTheSpaceUsingTheWebdavApi(
		string $user,
		string $source,
		string $destination,
		string $mtime,
		string $spaceName
	): void {
		$this->spacesContext->setSpaceIDByName($user, $spaceName);
		$this->tusContext->userUploadsFileWithContentToWithMtimeUsingTUS($user, $source, $destination, $mtime);
	}

	/**
	 * @Given /^user "([^"]*)" has uploaded file with checksum "([^"]*)" to the last created TUS Location with offset "([^"]*)" and content "([^"]*)" via TUS inside of the space "([^"]*)" using the WebDAV API$/
	 *
	 * @param string $user
	 * @param string $checksum
	 * @param string $offset
	 * @param string $content
	 * @param string $spaceName
	 *
	 * @return void
	 * @throws Exception|GuzzleException
	 */
	public function userHasUploadedFileWithChecksumToTheLastCreatedTusLocationWithOffsetAndContentViaTusInsideOfTheSpaceUsingTheWebdavApi(
		string $user,
		string $checksum,
		string $offset,
		string $content,
		string $spaceName
	): void {
		$this->spacesContext->setSpaceIDByName($user, $spaceName);
		$this->tusContext->userHasUploadedFileWithChecksum($user, $checksum, $offset, $content);
	}

	/**
	 * @Given /^user "([^"]*)" uploads file with checksum "([^"]*)" to the last created TUS Location with offset "([^"]*)" and content "([^"]*)" via TUS inside of the space "([^"]*)" using the WebDAV API$/
	 *
	 * @param string $user
	 * @param string $checksum
	 * @param string $offset
	 * @param string $content
	 * @param string $spaceName
	 *
	 * @return void
	 * @throws Exception|GuzzleException
	 */
	public function userUploadsFileWithChecksumToTheLastCreatedTusLocationWithOffsetAndContentViaTusInsideOfTheSpaceUsingTheWebdavApi(
		string $user,
		string $checksum,
		string $offset,
		string $content,
		string $spaceName
	): void {
		$this->spacesContext->setSpaceIDByName($user, $spaceName);
		$this->tusContext->userUploadsFileWithChecksum($user, $checksum, $offset, $content);
	}

	/**
	 * @When /^user "([^"]*)" sends a chunk to the last created TUS Location with offset "([^"]*)" and data "([^"]*)" with checksum "([^"]*)" via TUS inside of the space "([^"]*)" using the WebDAV API$/
	 *
	 * @param string $user
	 * @param string $offset
	 * @param string $data
	 * @param string $checksum
	 * @param string $spaceName
	 *
	 * @return void
	 * @throws Exception|GuzzleException
	 */
	public function userSendsAChunkToTheLastCreatedTusLocationWithOffsetAndDataWithChecksumViaTusInsideOfTheSpaceUsingTheWebdavApi(
		string $user,
		string $offset,
		string $data,
		string $checksum,
		string $spaceName
	): void {
		$this->spacesContext->setSpaceIDByName($user, $spaceName);
		$this->tusContext->userUploadsChunkFileWithChecksum($user, $offset, $data, $checksum);
	}

	/**
	 * @When /^user "([^"]*)" overwrites recently shared file with offset "([^"]*)" and data "([^"]*)" with checksum "([^"]*)" via TUS inside of the space "([^"]*)" using the WebDAV API with these headers:$/
	 *
	 * @param string $user
	 * @param string $offset
	 * @param string $data
	 * @param string $checksum
	 * @param string $spaceName
	 * @param TableNode $headers
	 *
	 * @return void
	 * @throws GuzzleException
	 */
	public function userOverwritesRecentlySharedFileWithOffsetAndDataWithChecksumViaTusInsideOfTheSpaceUsingTheWebdavApiWithTheseHeaders(
		string $user,
		string $offset,
		string $data,
		string $checksum,
		string $spaceName,
		TableNode $headers
	): void {
		$this->spacesContext->setSpaceIDByName($user, $spaceName);
		$this->tusContext->userOverwritesFileWithChecksum($user, $offset, $data, $checksum, $headers);
	}

	/**
	 * @Then /^as "([^"]*)" the mtime of the file "([^"]*)" in space "([^"]*)" should be "([^"]*)"$/
	 *
	 * @param string $user
	 * @param string $resource
	 * @param string $spaceName
	 * @param string $mtime
	 *
	 * @return void
	 * @throws Exception|GuzzleException
	 */
	public function theMtimeOfTheFileInSpaceShouldBe(
		string $user,
		string $resource,
		string $spaceName,
		string $mtime
	): void {
		$this->spacesContext->setSpaceIDByName($user, $spaceName);
		$this->featureContext->theMtimeOfTheFileShouldBe($user, $resource, $mtime);
	}
}

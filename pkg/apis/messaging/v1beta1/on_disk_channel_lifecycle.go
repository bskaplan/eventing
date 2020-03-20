/*
 * Copyright 2020 The Knative Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

 package v1beta1

 import (
	 appsv1 "k8s.io/api/apps/v1"
	 corev1 "k8s.io/api/core/v1"
	 "k8s.io/apimachinery/pkg/runtime/schema"
	 "knative.dev/pkg/apis"
	 v1 "knative.dev/pkg/apis/duck/v1"
 )
 
 var odcCondSet = apis.NewLivingConditionSet(OnDiskChannelConditionDispatcherReady, OnDiskChannelConditionServiceReady, OnDiskChannelConditionEndpointsReady, OnDiskChannelConditionAddressable, OnDiskChannelConditionChannelServiceReady)
 
 const (
	 // OnDiskChannelConditionReady has status True when all subconditions below have been set to True.
	 OnDiskChannelConditionReady = apis.ConditionReady
 
	 // OnDiskChannelConditionDispatcherReady has status True when a Dispatcher deployment is ready
	 // Keyed off appsv1.DeploymentAvaialble, which means minimum available replicas required are up
	 // and running for at least minReadySeconds.
	 OnDiskChannelConditionDispatcherReady apis.ConditionType = "DispatcherReady"
 
	 // OnDiskChannelConditionServiceReady has status True when a k8s Service is ready. This
	 // basically just means it exists because there's no meaningful status in Service. See Endpoints
	 // below.
	 OnDiskChannelConditionServiceReady apis.ConditionType = "ServiceReady"
 
	 // OnDiskChannelConditionEndpointsReady has status True when a k8s Service Endpoints are backed
	 // by at least one endpoint.
	 OnDiskChannelConditionEndpointsReady apis.ConditionType = "EndpointsReady"
 
	 // OnDiskChannelConditionAddressable has status true when this OnDiskChannel meets
	 // the Addressable contract and has a non-empty hostname.
	 OnDiskChannelConditionAddressable apis.ConditionType = "Addressable"
 
	 // OnDiskChannelConditionServiceReady has status True when a k8s Service representing the channel is ready.
	 // Because this uses ExternalName, there are no endpoints to check.
	 OnDiskChannelConditionChannelServiceReady apis.ConditionType = "ChannelServiceReady"
 )
 
 // GetGroupVersionKind returns GroupVersionKind for OnDiskChannels
 func (odc *OnDiskChannel) GetGroupVersionKind() schema.GroupVersionKind {
	 return SchemeGroupVersion.WithKind("OnDiskChannel")
 }
 
 // GetUntypedSpec returns the spec of the OnDiskChannel.
 func (i *OnDiskChannel) GetUntypedSpec() interface{} {
	 return i.Spec
 }
 
 // GetCondition returns the condition currently associated with the given type, or nil.
 func (odcs *OnDiskChannelStatus) GetCondition(t apis.ConditionType) *apis.Condition {
	 return odcCondSet.Manage(odcs).GetCondition(t)
 }
 
 // IsReady returns true if the resource is ready overall.
 func (odcs *OnDiskChannelStatus) IsReady() bool {
	 return odcCondSet.Manage(odcs).IsHappy()
 }
 
 // InitializeConditions sets relevant unset conditions to Unknown state.
 func (odcs *OnDiskChannelStatus) InitializeConditions() {
	 odcCondSet.Manage(odcs).InitializeConditions()
 }
 
 func (odcs *OnDiskChannelStatus) SetAddress(url *apis.URL) {
	 odcs.Address = &v1.Addressable{URL: url}
	 if url != nil {
		 odcCondSet.Manage(odcs).MarkTrue(OnDiskChannelConditionAddressable)
	 } else {
		 odcCondSet.Manage(odcs).MarkFalse(OnDiskChannelConditionAddressable, "emptyHostname", "hostname is the empty string")
	 }
 }
 
 func (odcs *OnDiskChannelStatus) MarkDispatcherFailed(reason, messageFormat string, messageA ...interface{}) {
	 odcCondSet.Manage(odcs).MarkFalse(OnDiskChannelConditionDispatcherReady, reason, messageFormat, messageA...)
 }
 
 func (odcs *OnDiskChannelStatus) MarkDispatcherUnknown(reason, messageFormat string, messageA ...interface{}) {
	 odcCondSet.Manage(odcs).MarkUnknown(OnDiskChannelConditionDispatcherReady, reason, messageFormat, messageA...)
 }
 
 // TODO: Unify this with the ones from Eventing. Say: Broker, Trigger.
 func (odcs *OnDiskChannelStatus) PropagateDispatcherStatus(ds *appsv1.DeploymentStatus) {
	 for _, cond := range ds.Conditions {
		 if cond.Type == "Available" {
			 if cond.Status == corev1.ConditionTrue {
				 odcCondSet.Manage(odcs).MarkTrue(OnDiskChannelConditionDispatcherReady)
			 } else if cond.Status == corev1.ConditionFalse {
				 odcs.MarkDispatcherFailed("DispatcherDeploymentFalse", "The status of Dispatcher Deployment is False: %s : %s", cond.Reason, cond.Message)
			 } else if cond.Status == corev1.ConditionUnknown {
				 odcs.MarkDispatcherUnknown("DispatcherDeploymentUnknown", "The status of Dispatcher Deployment is Unknown: %s : %s", cond.Reason, cond.Message)
			 }
		 }
	 }
 }
 
 func (odcs *OnDiskChannelStatus) MarkServiceFailed(reason, messageFormat string, messageA ...interface{}) {
	 odcCondSet.Manage(odcs).MarkFalse(OnDiskChannelConditionServiceReady, reason, messageFormat, messageA...)
 }
 
 func (odcs *OnDiskChannelStatus) MarkServiceUnknown(reason, messageFormat string, messageA ...interface{}) {
	 odcCondSet.Manage(odcs).MarkUnknown(OnDiskChannelConditionServiceReady, reason, messageFormat, messageA...)
 }
 
 func (odcs *OnDiskChannelStatus) MarkServiceTrue() {
	 odcCondSet.Manage(odcs).MarkTrue(OnDiskChannelConditionServiceReady)
 }
 
 func (odcs *OnDiskChannelStatus) MarkChannelServiceFailed(reason, messageFormat string, messageA ...interface{}) {
	 odcCondSet.Manage(odcs).MarkFalse(OnDiskChannelConditionChannelServiceReady, reason, messageFormat, messageA...)
 }
 
 func (odcs *OnDiskChannelStatus) MarkChannelServiceUnknown(reason, messageFormat string, messageA ...interface{}) {
	 odcCondSet.Manage(odcs).MarkUnknown(OnDiskChannelConditionChannelServiceReady, reason, messageFormat, messageA...)
 }
 
 func (odcs *OnDiskChannelStatus) MarkChannelServiceTrue() {
	 odcCondSet.Manage(odcs).MarkTrue(OnDiskChannelConditionChannelServiceReady)
 }
 
 func (odcs *OnDiskChannelStatus) MarkEndpointsFailed(reason, messageFormat string, messageA ...interface{}) {
	 odcCondSet.Manage(odcs).MarkFalse(OnDiskChannelConditionEndpointsReady, reason, messageFormat, messageA...)
 }
 
 func (odcs *OnDiskChannelStatus) MarkEndpointsUnknown(reason, messageFormat string, messageA ...interface{}) {
	 odcCondSet.Manage(odcs).MarkUnknown(OnDiskChannelConditionEndpointsReady, reason, messageFormat, messageA...)
 }
 
 func (odcs *OnDiskChannelStatus) MarkEndpointsTrue() {
	 odcCondSet.Manage(odcs).MarkTrue(OnDiskChannelConditionEndpointsReady)
 }
 
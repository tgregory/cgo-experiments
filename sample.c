#define INTERNAL
#include "sample.h"
#include "stdlib.h"

struct Sample {
	int magic_number;
	SampleCallback callback;
	void* callback_data;
};

int dummy_callback(Sample* sample, int number, void* data){
	return 0;
}

Sample* create_sample(int number) {
	Sample* result = (Sample*)malloc(sizeof(Sample));
	result->magic_number=number;
	result->callback = dummy_callback;
	return result;
}

int destroy_sample(Sample* sample) {
	free(sample);
	return 0;
}

int invoke_callback(Sample* sample) {
	return sample->callback(sample, (sample->magic_number)++, sample->callback_data);
}

int register_callback(Sample* sample, SampleCallback callback, void* arbitrary_data){
	sample->callback = callback;
	sample->callback_data = arbitrary_data;
	return 0;
}
